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
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
)

var (
	ginRoute *gin.Engine
	conf     = &config.Config{}
)

func main() {
	err := conf.LoadEnvDefault()
	if err != nil {
		log.Fatal(err)
	}
	// 初始化日志格式
	config.InitLogrus()
	logrus.Info("config loaded:", conf)

	sql.Init(conf)
	redis.Init(conf)
	etcd.Init(conf)
	err = judge.Init()
	if err == nil {
		go judge.Start()
	} else {
		logrus.Error(err, "评测机器初始化失败")
	}
	ginRoute = gin.Default()
	controller.Init(ginRoute)
	logrus.Info(fmt.Sprintf(":%s", conf.BackendServerPort))
	ginRoute.Run(fmt.Sprintf(":%s", conf.BackendServerPort))
}
