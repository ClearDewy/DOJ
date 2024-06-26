/**
 * @ Author: ClearDewy
 * @ Desc: 初始化操作
 **/
package controller

import (
	"context"
	"doj-go/DataBackup/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Init() (func() error, func(ctx context.Context) error) {
	r := gin.Default()
	r.Static("/api/public", "./doj/public")

	// 普通请求模块
	commentGroup := r.Group("/api")
	{
		commentGroup.GET("/get-web-config", getWebConfig)
	}
	// 用户请求模块
	userGroup := r.Group("/api")
	{
		userGroup.POST("/login", login)
		userGroup.GET("/logout", requireAuth, logout)
	}
	// 题目请求模块
	problemGroup := r.Group("/api")
	{
		problemGroup.GET("/get-problem-list", getProblemList)
		problemGroup.GET("/get-problem-detail", requireAuth, getProblemDetail)
	}
	// 判题请求模块
	judgeGroup := r.Group("/api", requireAuth)
	{
		judgeGroup.POST("/submit-problem-judge", submitProblemJudge)
	}
	return func() error {
		return r.Run(fmt.Sprintf(":%s", config.Conf.BackendServerPort))
	}, nil
}
