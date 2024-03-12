/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package controller

import (
	"doj-go/DataBackup/judge"
	"doj-go/DataBackup/sql"
	"doj-go/DataBackup/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubmitProblemJudgeType struct {
	ProblemId string `json:"problem_id"`
	Code      string `json:"code"`
	Lid       int    `json:"lid"`
}

func submitProblemJudge(c *gin.Context) {
	spj := &SubmitProblemJudgeType{}
	err := c.BindJSON(spj)
	if err != nil {
		utils.HandleError(err, "")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "提交代码失败",
		})
		return
	}
	value, _ := c.Get(USERINFO)
	user_info := value.(sql.UserInfoType)
	pj, err := sql.GetProblemJudgeByProblemId(spj.ProblemId)
	if err != nil {
		utils.HandleError(err, "")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "提交代码失败",
		})
		return
	}

	jid, err := sql.AddJudge(&sql.AddJudgeType{
		Uid:       string(user_info.Uid),
		Username:  string(user_info.Username),
		Pid:       pj.Pid,
		ProblemId: spj.ProblemId,
		Code:      spj.Code,
		Lid:       spj.Lid,
	})
	if err != nil {
		utils.HandleError(err, "")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "提交代码失败",
		})
		return
	}

	judge.AddCommonProblemJudge(&judge.JudgeItemType{
		Uid:         string(user_info.Uid),
		Pid:         pj.Pid,
		Jid:         int(jid),
		Parallelism: pj.Parallelism,
	})

	c.JSON(http.StatusOK, gin.H{
		"jid": jid,
	})
}
