/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package main

import (
	"doj-go/DataBackup/config"
	"doj-go/DataBackup/controller"
	"doj-go/DataBackup/etcd"
	"doj-go/DataBackup/judge"
	"doj-go/DataBackup/redis"
	"doj-go/DataBackup/sql"
	"fmt"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	ginRoute *gin.Engine
)

func main() {
	err := config.Conf.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	// 初始化日志格式
	logrus.Info("config loaded:", config.Conf)

	sql.Init(config.Conf)
	redis.Init(config.Conf)
	etcd.Init(config.Conf)
	err = judge.Init()
	if err == nil {
		go judge.Start()
	} else {
		logrus.Error(err, "评测机器初始化失败")
	}
	ginRoute = gin.Default()
	controller.Init(ginRoute)
	logrus.Info(fmt.Sprintf(":%s", config.Conf.BackendServerPort))
	ginRoute.Run(fmt.Sprintf(":%s", config.Conf.BackendServerPort))
}
