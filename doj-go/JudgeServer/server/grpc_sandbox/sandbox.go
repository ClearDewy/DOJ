/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_sandbox

import (
	"context"
	"doj-go/JudgeServer/pb"
	"doj-go/JudgeServer/server/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	//SandBoxGrpcAddr = "127.0.0.1:5051"
	SandBoxGrpcAddr = "172.24.0.7:5051"
	MaxProcLimit    = 128
)

var (
	Client pb.ExecutorClient
)

func Init() {
	conn, err := grpc.Dial(SandBoxGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.HandleError(err, "sandbox连接失败")
	}
	Client = pb.NewExecutorClient(conn)
}

func DeleteFileByFileID(fileIDs []string) {
	for _, value := range fileIDs {
		_, err := Client.FileDelete(context.Background(), &pb.FileID{
			FileID: value,
		})
		if err != nil {
			utils.HandleError(err, "删除缓存文件异常")
		}
	}
}
func DeleteFileByResult(result *pb.Response_Result) {
	if result == nil {
		return
	}
	for _, value := range result.FileIDs {
		_, err := Client.FileDelete(context.Background(), &pb.FileID{
			FileID: value,
		})
		if err != nil {
			utils.HandleError(err, "删除缓存文件异常")
		}
	}
}
