/**
 * @ Author: ClearDewy
 * @ Desc: 判题队列
 **/
package judge

import (
	"context"
	"doj-go/DataBackup/etcd"
	"doj-go/DataBackup/redis"
	"doj-go/DataBackup/utils"
	"doj-go/jpb"
	"encoding/json"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"strings"
	"sync"
)

const (
	COMMON_PROBLEM_JUDGE  = "common_problem_judge"
	CONTEST_PROBLEM_JUDGE = "contest_problem_judge"
	JUDGE_SERVER          = "judge_server"
	JUDGE_SERVER_PRE      = "/judge/judge-server/"
	JUDGE_PARALLELISM_PRE = "/judge/parallelism/"
)

type JudgeItemType struct {
	Uid         string `json:"uid"`
	Jid         int    `json:"jid"`
	Pid         int    `json:"pid"`
	Parallelism int    `json:"parallelism"`
}

type JudgeServerType struct {
	Client      jpb.JudgeServerClient
	Parallelism int
	Error       bool
	Mutex       sync.Mutex
}

var judgeServers = make(map[string]*JudgeServerType)

func Init() error {
	resp, err := etcd.Client.Get(context.Background(), JUDGE_SERVER_PRE, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		addr := strings.Replace(string(kv.Key), JUDGE_SERVER_PRE, "", 1)

		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			utils.HandleError(err, "sandbox连接失败")
			continue
		}

		judgeServers[addr] = &JudgeServerType{
			Client:      jpb.NewJudgeServerClient(conn),
			Parallelism: GetParallelism(addr),
		}

		redis.Rdb.RPush(context.Background(), JUDGE_SERVER, addr)
	}
	go func() {
		watchChan := etcd.Client.Watch(context.Background(), JUDGE_SERVER_PRE, clientv3.WithPrefix())
		for wresp := range watchChan {
			for _, ev := range wresp.Events {
				addr := strings.Replace(string(ev.Kv.Key), JUDGE_SERVER_PRE, "", 1)
				// 新增
				if ev.Type == clientv3.EventTypePut {
					conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						utils.HandleError(err, "sandbox连接失败")
						continue
					}
					judgeServers[addr] = &JudgeServerType{
						Client:      jpb.NewJudgeServerClient(conn),
						Parallelism: GetParallelism(addr),
					}
					redis.Rdb.RPush(context.Background(), JUDGE_SERVER, addr)
				} else {
					// 删除
					redis.Rdb.LRem(context.Background(), JUDGE_SERVER, 1, addr)
					delete(judgeServers, addr)
				}
			}
		}
	}()
	return nil
}

func GetParallelism(addr string) int {
	resp, err := etcd.Client.Get(context.Background(), JUDGE_PARALLELISM_PRE+addr)
	if err != nil {
		utils.HandleError(err, "获取判题机信息失败")
	}
	if len(resp.Kvs) == 0 {
		utils.HandleError(err, "获取判题机信息失败")
		return 0
	}
	parallelism, err := strconv.Atoi(string(resp.Kvs[0].Value))
	if err != nil {
		utils.HandleError(err, "获取判题机信息失败")
	}
	return parallelism
}

func Start() {
	for {
		// 优先评测比赛的提交
		// 如果没有提交则会阻塞
		// resp[0] 为key，resp[1] 为value
		resp, err := redis.Rdb.BLPop(context.Background(), 0, CONTEST_PROBLEM_JUDGE, COMMON_PROBLEM_JUDGE).Result()
		if err != nil {
			utils.HandleError(err, "获取评测队列元素异常")
			continue
		}
		logrus.Info(resp)
		judgeItem := &JudgeItemType{}
		jsonItem := resp[1]
		err = json.Unmarshal([]byte(resp[1]), judgeItem)
		if err != nil {
			utils.HandleError(err, "解析评测队列元素异常")
			continue
		}
		resp, err = redis.Rdb.BLPop(context.Background(), 0, JUDGE_SERVER).Result()
		if err != nil {
			utils.HandleError(err, "获取评测队列元素异常")
			continue
		}
		logrus.Info(resp)
		go Judge(judgeItem, resp[1], resp[0], jsonItem)
	}
}

func AddCommonProblemJudge(item *JudgeItemType) {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		utils.HandleError(err, "对象JSON化失败")
		return
	}
	push := redis.Rdb.RPush(context.Background(), COMMON_PROBLEM_JUDGE, jsonItem)
	logrus.Info(push)
}
func AddContestProblemJudge(item *jpb.JudgeItem) {
	jsonItem, err := json.Marshal(item)
	if err != nil {
		utils.HandleError(err, "对象JSON化失败")
		return
	}
	redis.Rdb.RPush(context.Background(), CONTEST_PROBLEM_JUDGE, jsonItem)
}

func Judge(judgeItem *JudgeItemType, addr string, problem_key string, jsonItem string) {
	server := judgeServers[addr]
	// 占用评测资源
	server.Mutex.Lock()
	server.Parallelism -= judgeItem.Parallelism
	server.Mutex.Unlock()
	if server.Parallelism > 0 {
		// 如果还有评测资源，则加入redis判题服务队列
		redis.Rdb.RPush(context.Background(), JUDGE_SERVER, addr)
	}
	_, err := server.Client.Judge(context.Background(), &jpb.JudgeItem{
		Uid: judgeItem.Uid,
		Pid: int32(judgeItem.Pid),
		Jid: int32(judgeItem.Jid),
	})
	if err != nil {
		// 如果评测出现错误，则重新加入判题队列
		redis.Rdb.RPush(context.Background(), problem_key, jsonItem)
		utils.HandleError(err, "判题机出现异常")
		server.Error = true
		return
	}
	// 顺利评测完,归还评测资源
	server.Mutex.Lock()
	server.Parallelism -= judgeItem.Parallelism
	server.Mutex.Unlock()
}
