/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package controller

import (
	"context"
	"doj-go/DataBackup/config"
	"doj-go/DataBackup/etcd"
	"doj-go/DataBackup/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getWebConfig(c *gin.Context) {
	webConfig := &config.WebConfig{}
	resp, err := etcd.Client.Get(context.Background(), etcd.WEB_CONFIG)
	if err != nil {
		utils.HandleError(err, "")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err,
		})
		return
	}
	if len(resp.Kvs) != 1 {
		utils.HandleError(nil, "服务器配置不唯一")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器配置不唯一，请联系管理员",
		})
		return
	}
	err = json.Unmarshal(resp.Kvs[0].Value, webConfig)
	if err != nil {
		utils.HandleError(err, "")
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err,
		})
		return
	}
	c.JSON(http.StatusOK, webConfig)
}
