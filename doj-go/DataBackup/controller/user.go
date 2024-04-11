/**
 * @ Author: ClearDewy
 * @ Desc: 用户的api
 **/
package controller

import (
	"context"
	"doj-go/DataBackup/redis"
	"doj-go/DataBackup/sql"
	"doj-go/DataBackup/utils"
	"encoding/json"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
	"time"
)

const (
	TRY_LOGIN_TIME = "try_login_time:"

	AUTHORIZATION = "authorization"
	USERINFO      = "user_info"
)

type LoginType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserInfoType struct {
	Username            string `json:"username"`
	School              string `json:"school"`
	Major               string `json:"major"`
	Number              string `json:"number"`
	Name                string `json:"name"`
	Gender              string `json:"gender"`
	Cf_username         string `json:"cf_username"`
	Email               string `json:"email"`
	Avatar              string `json:"avatar"`
	Signature           string `json:"signature"`
	Title_name          string `json:"title_name"`
	Title_color         string `json:"title_color"`
	System_auth         int    `json:"system_auth"`
	User_auth           int    `json:"user_auth"`
	Problem_auth        int    `json:"problem_auth"`
	Context_auth        int    `json:"context_auth"`
	Train_auth          int    `json:"train_auth"`
	Submit_auth         int    `json:"submit_auth"`
	Context_attend_auth int    `json:"context_attend_auth"`
	Train_attend_auth   int    `json:"train_attend_auth"`
}

func login(c *gin.Context) {
	var loginDto = LoginType{}
	if err := c.BindJSON(&loginDto); err != nil {
		logrus.ErrorM(err, "参数校验失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if len(loginDto.Password) < 6 || len(loginDto.Password) > 20 || len(loginDto.Username) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户名长度不得超过20，密码长度不得超过20、少于6",
		})
		return
	}
	try_login_time_key := TRY_LOGIN_TIME + c.RemoteIP()
	try_login_time := 0

	if value, err := redis.Rdb.Get(context.Background(), try_login_time_key).Result(); err == nil {
		try_login_time, _ = strconv.Atoi(value)
	}
	if try_login_time > 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "对不起！登录失败次数过多！您的账号有风险，半个小时内暂时无法登录！",
		})
		return
	}
	userInfoDao, err := sql.GetUserInfoByUsername(loginDto.Username)
	if err != nil {
		logrus.ErrorM(err, "")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}
	if string(userInfoDao.Password) == loginDto.Password {
		token, err := utils.JwtGenerate(string(userInfoDao.Uid))
		if err != nil {
			logrus.ErrorM(err, "")
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "生成 token 失败",
			})
			return
		}
		c.Header(AUTHORIZATION, token)
		var userInfoDto = UserInfoType{}
		userInfoByte, _ := json.Marshal(userInfoDao)
		_ = json.Unmarshal(userInfoByte, &userInfoDto)
		c.JSON(http.StatusOK, userInfoDto)
		redis.Rdb.Del(context.Background(), try_login_time_key)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户名或密码错误",
		})
		redis.Rdb.Set(context.Background(), try_login_time_key, strconv.Itoa(try_login_time+1), 30*time.Minute)
	}
}

func logout(c *gin.Context) {
	c.Header("authorization", "")

	c.JSON(http.StatusOK, nil)
}
