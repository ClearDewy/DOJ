/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package controller

import (
	"doj-go/DataBackup/sql"
	"doj-go/DataBackup/utils"
	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

// 验证登录的token
func requireAuth(c *gin.Context) {
	token := c.GetHeader(AUTHORIZATION)
	jwtParse, err := utils.JwtParse(token)
	if err != nil {
		utils.HandleError(err, "")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "token解析失败",
		})

		return
	}
	if jwtParse.ExpiresAt.Before(time.Now()) {
		// token过期了，但是未达到过期时间一倍，则自动续约
		if jwtParse.ExpiresAt.Before(time.Now().Add(utils.ExpireTime)) {
			newToken, err := utils.JwtGenerate(jwtParse.Uid)
			if err != nil {
				utils.HandleError(err, "")
			}
			c.Header(AUTHORIZATION, newToken)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "登录已过期，请重新登录",
			})
			return
		}
	}
	userInfo, err := sql.GetUserInfoByUid(jwtParse.Uid)
	if err != nil {
		utils.HandleError(err, "")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "用户不存在，请重新登录",
		})
		return
	}
	c.Set(USERINFO, userInfo)
	c.Next()
}

// 尝试验证登录的token
func tryAuth(c *gin.Context) {
	token := c.GetHeader(AUTHORIZATION)
	jwtParse, err := utils.JwtParse(token)
	if err != nil {
		c.Next()
		return
	}
	if jwtParse.ExpiresAt.Before(time.Now()) {
		// token过期了，但是未达到过期时间一倍，则自动续约
		if jwtParse.ExpiresAt.Before(time.Now().Add(utils.ExpireTime)) {
			newToken, err := utils.JwtGenerate(jwtParse.Uid)
			if err != nil {
				utils.HandleError(err, "")
			}
			c.Header(AUTHORIZATION, newToken)
		} else {
			c.Header(AUTHORIZATION, "")
			c.Next()
			return
		}
	}
	userInfo, err := sql.GetUserInfoByUid(jwtParse.Uid)
	if err != nil {
		c.Header(AUTHORIZATION, "")
		c.Next()
		return
	}
	c.Set(USERINFO, userInfo)
	c.Next()
}
