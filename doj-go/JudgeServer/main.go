/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package main

import (
	"context"
	"doj-go/JudgeServer/config"
	"doj-go/JudgeServer/grpc_judge"
	"doj-go/JudgeServer/internal/etcd"
	"doj-go/JudgeServer/internal/redis"
	"doj-go/JudgeServer/internal/sql"
	"doj-go/JudgeServer/sandbox"
	"doj-go/jspb"
	"github.com/ClearDewy/go-pkg/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := config.Conf.LoadEnv()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Config:", config.Conf)
	servers := []func() (func() error, func(ctx context.Context) error){
		redis.Init,
		sql.Init,
		etcd.Init,
		sandbox.Init,
		grpc_judge.Init,
	}

	sig := make(chan os.Signal, len(servers)+1)
	stops := make([]func(ctx context.Context) error, 0, len(servers))

	for _, init := range servers {
		start, stop := init()
		if start != nil {
			go func() {
				err := start()
				logrus.ErrorM(err, "")
				sig <- os.Interrupt
			}()
		}
		if stop != nil {
			stops = append(stops, stop)
		}
	}

	ju := grpc_judge.JudgeServer{}
	ju.Judge(context.Background(), &jspb.JudgeItem{
		Uid: "1",
		Jid: 2,
		Pid: 1000,
	})

	// Graceful shutdown...
	signal.Notify(sig, os.Interrupt)
	<-sig
	signal.Reset(os.Interrupt)

	logrus.Info("判题服务关闭中……")

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	var eg errgroup.Group
	for _, s := range stops {
		eg.Go(func() error {
			return s(ctx)
		})
	}

	go func() {
		logrus.Info("判题服务关闭完成", eg.Wait())
		cancel()
	}()
	<-ctx.Done()
}
