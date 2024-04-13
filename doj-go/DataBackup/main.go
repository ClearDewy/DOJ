/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package main

import (
	"context"
	"doj-go/DataBackup/config"
	"doj-go/DataBackup/controller"
	"doj-go/DataBackup/internal/etcd"
	"doj-go/DataBackup/internal/redis"
	"doj-go/DataBackup/internal/sql"
	"doj-go/DataBackup/judge"
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
		judge.Init,
		controller.Init,
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

	// Graceful shutdown...
	signal.Notify(sig, os.Interrupt)
	<-sig
	signal.Reset(os.Interrupt)

	logrus.Info("后台服务关闭中……")

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	var eg errgroup.Group
	for _, s := range stops {
		eg.Go(func() error {
			return s(ctx)
		})
	}

	go func() {
		logrus.Info("后台服务关闭完成", eg.Wait())
		cancel()
	}()
	<-ctx.Done()
}
