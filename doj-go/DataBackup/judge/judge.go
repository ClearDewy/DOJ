/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package judge

import (
	"context"
	"doj-go/DataBackup/internal/etcd"
	"doj-go/DataBackup/internal/redis"
	"github.com/ClearDewy/go-pkg/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

func Init() (func() error, func(ctx context.Context) error) {
	go CloseExpireConn()
	err := redis.Rdb.Del(context.Background(), JUDGE_SERVER).Err()
	if err != nil {
		logrus.FatalM(err, "清空judge server队列失败")
	}
	resp, err := etcd.Client.Get(context.Background(), JUDGE_SERVER_PRE, clientv3.WithPrefix())

	if err != nil {
		logrus.FatalM(err, "init judge server connection failed")
	}
	for _, kv := range resp.Kvs {
		addr := strings.Replace(string(kv.Key), JUDGE_SERVER_PRE, "", 1)

		ConnectJudgeServer(addr)
	}

	return Start, Stop
}

func Start() error {
	watchChan := etcd.Client.Watch(context.Background(), JUDGE_SERVER_PRE, clientv3.WithPrefix())
	for wresp := range watchChan {
		for _, ev := range wresp.Events {
			addr := strings.Replace(string(ev.Kv.Key), JUDGE_SERVER_PRE, "", 1)
			// 新增
			if ev.Type == clientv3.EventTypePut {
				ConnectJudgeServer(addr)
			} else {
				// 删除
				redis.Rdb.LRem(context.Background(), JUDGE_SERVER, 1, addr)
				judgeServerConnPoolMutex.Lock()
				delete(judgeServerConnPool, addr)
				judgeServerConnPoolMutex.Unlock()
			}
		}
	}
	return nil
}

func Stop(ctx context.Context) error {
	judgeServerConnPoolMutex.Lock()
	defer judgeServerConnPoolMutex.Unlock()
	for _, v := range judgeServerConnPool {
		v.Conn.Close()
	}
	return nil
}
