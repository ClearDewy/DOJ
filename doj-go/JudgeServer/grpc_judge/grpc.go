/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"context"
	"doj-go/JudgeServer/config"
	"doj-go/JudgeServer/internal/redis"
	"doj-go/jspb"
	"google.golang.org/grpc"
	"net"
)

const (
	JUDGE_SERVER = "judge_server"
)

func Init() (func() error, func(ctx context.Context) error) {
	grpcServer := grpc.NewServer()

	return func() error {
			// 讲判题机加入进可用队列中
			Useable = true
			redis.Rdb.RPush(context.Background(), JUDGE_SERVER, config.Conf.JudgeServerAddr)

			listen, err := net.Listen("tcp", config.Conf.JudgeServerAddr)
			if err != nil {
				return err
			}
			jspb.RegisterJudgeServerServer(grpcServer, &JudgeServer{})
			err = grpcServer.Serve(listen)
			return err
		}, func(ctx context.Context) error {
			grpcServer.Stop()
			return nil
		}
}
