/**
 * @ Author: ClearDewy
 * @ Desc: 判题队列
 **/
package judge

import (
	"context"
	"doj-go/DataBackup/internal/etcd"
	"doj-go/DataBackup/internal/redis"
	"doj-go/jspb"
	"encoding/json"
	"github.com/ClearDewy/go-pkg/logrus"
	"strconv"
)

const (
	COMMON_PROBLEM_JUDGE  = "common_problem_judge"
	CONTEST_PROBLEM_JUDGE = "contest_problem_judge"
	JUDGE_SERVER          = "judge_server"
	JUDGE_SERVER_PRE      = "/judge/judge-server/"
	JUDGE_PARALLELISM_PRE = "/judge/parallelism/"
)

func GetParallelism(addr string) int {
	resp, err := etcd.Client.Get(context.Background(), JUDGE_PARALLELISM_PRE+addr)
	if err != nil {
		logrus.ErrorM(err, "获取判题机信息失败")
	}
	if len(resp.Kvs) == 0 {
		logrus.ErrorM(err, "获取判题机信息失败")
		return 0
	}
	parallelism, err := strconv.Atoi(string(resp.Kvs[0].Value))
	if err != nil {
		logrus.ErrorM(err, "获取判题机信息失败")
	}
	return parallelism
}

func Run() {
	for {
		// 优先评测比赛的提交
		// 如果没有提交则会阻塞
		// resp[0] 为key，resp[1] 为value
		pResp, err := redis.Rdb.BLPop(context.Background(), 0, CONTEST_PROBLEM_JUDGE, COMMON_PROBLEM_JUDGE).Result()
		if err != nil {
			logrus.ErrorM(err, "获取评测队列元素异常")
			continue
		}
		logrus.Info(pResp)
		judgeItem := &jspb.JudgeItem{}
		err = json.Unmarshal([]byte(pResp[1]), judgeItem)
		if err != nil {
			logrus.ErrorM(err, "解析评测队列元素异常")
			continue
		}
		jsResp, err := redis.Rdb.BRPop(context.Background(), 0, JUDGE_SERVER).Result()
		if err != nil {
			logrus.ErrorM(err, "获取评测队列元素异常")
			continue
		}
		logrus.Info(pResp)
		go Judge(judgeItem, jsResp[1], pResp[0], pResp[1])
	}
}

func AddCommonProblemJudge(item *jspb.JudgeItem) {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		logrus.ErrorM(err, "对象JSON化失败")
		return
	}
	push := redis.Rdb.RPush(context.Background(), COMMON_PROBLEM_JUDGE, jsonItem)
	logrus.Info(push)
}

func AddContestProblemJudge(item *jspb.JudgeItem) {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		logrus.ErrorM(err, "对象JSON化失败")
		return
	}
	redis.Rdb.RPush(context.Background(), CONTEST_PROBLEM_JUDGE, jsonItem)
}

func Judge(judgeItem *jspb.JudgeItem, addr string, k, v string) {
	judgeServerConnPoolMutex.RLocker()
	server := judgeServerConnPool[addr]
	judgeServerConnPoolMutex.RUnlock()
	CheckJudgeServerConn(server)

	_, err := server.Client.Judge(context.Background(), judgeItem)
	if err != nil {
		// 如果评测出现错误
		logrus.ErrorM(err, "判题机出现异常;ip:"+addr)
		server.Error = true
		// 重新加入判题队列
		redis.Rdb.LPush(context.Background(), k, v)
		return
	}

}
