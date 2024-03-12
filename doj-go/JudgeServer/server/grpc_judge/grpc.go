/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"doj-go/JudgeServer/server/config"
	"doj-go/jpb"
	"google.golang.org/grpc"
	"net"
)

func Init() chan error {
	errChan := make(chan error)
	listen, err := net.Listen("tcp", config.Conf.JudgeServerAddr)
	if err != nil {
		errChan <- err
	}
	grpcServer := grpc.NewServer()
	jpb.RegisterJudgeServerServer(grpcServer, &JudgeServer{})
	if err := grpcServer.Serve(listen); err != nil {
		errChan <- err
	}
	return errChan
}
