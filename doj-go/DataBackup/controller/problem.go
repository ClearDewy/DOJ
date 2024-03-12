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
	"strconv"
)

func getProblemList(c *gin.Context) {
	oj := c.DefaultQuery("oj", "all")
	difficulty := c.DefaultQuery("difficulty", "all")
	keyword := c.DefaultQuery("keyword", "")
	tags := c.QueryArray("tags")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	currentPage, _ := strconv.Atoi(c.DefaultQuery("currentPage", "1"))

	problemList, err := sql.GetProblemList(oj, difficulty, keyword, tags, limit, currentPage)
	if err != nil {
		utils.HandleError(err, "")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取题目列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, problemList)
}

func getProblemDetail(c *gin.Context) {
	problem_id := c.Query("problem_id")
	pstmt, err := sql.GetProblemDetail(problem_id)
	if err != nil {
		utils.HandleError(err, "")
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取题目信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, pstmt)
}
