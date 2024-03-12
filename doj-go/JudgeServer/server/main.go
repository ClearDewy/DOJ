/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package main

import (
	"context"
	"doj-go/DataBackup/utils"
	"doj-go/JudgeServer/server/config"
	"doj-go/JudgeServer/server/etcd"
	"doj-go/JudgeServer/server/grpc_judge"
	"doj-go/JudgeServer/server/grpc_sandbox"
	"doj-go/JudgeServer/server/sql"
	"doj-go/jpb"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	err := config.Conf.LoadEnvDefault()
	if err != nil {
		log.Fatal(err)
	}
	// 初始化日志格式
	config.InitLogrus()
	logrus.Info("配置加载成功:", config.Conf)

	err = etcd.Init()
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("etcd连接成功")
	sql.Init()

	// 启动容器
	//errSandboxChan := grpc_sandbox.StartSandbox()
	grpc_sandbox.Init()

	logrus.Info("grpc sandbox连接成功")

	ju := grpc_judge.JudgeServer{}
	ju.Judge(context.Background(), &jpb.JudgeItem{
		Uid: "1",
		Jid: 2,
		Pid: 1000,
	})

	// 最后初始化这个，因为会阻塞
	errJudgeChan := grpc_judge.Init()

	select {
	//case err := <-errSandboxChan:
	//	utils.HandleError(err, "Sandbox Error")
	case err := <-errJudgeChan:
		utils.HandleError(err, "JudgeServer Error")
	}

	defer func() {
		etcd.Stop()
	}()
}
