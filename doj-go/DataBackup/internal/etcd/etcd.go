/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package etcd

import (
	"context"
	"doj-go/DataBackup/config"
	"encoding/json"
	"github.com/ClearDewy/go-pkg/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var Client *clientv3.Client

const (
	JUDGE_SERVER_PRE string = "/judge/judge-server/"
	CPU_NUM_PRE             = "/judge/cpu-num/"
	CPU_USAGE_PRE           = "/judge/cpu-usage/"
	MEM_USAGE_PRE           = "/judge/mem-usage/"
	WEB_CONFIG              = "/config/web"
	SERVICE_CONFIG          = "/config/service"
	JUDGE_CONFIG            = "/config/judge"
)

func Init() (func() error, func(ctx context.Context) error) {
	var err error
	Client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Conf.EtcdAddr},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    config.Conf.EtcdRootPassword,
	})
	if err != nil {
		logrus.FatalM(err, "etcd 连接失败")
	}
	initConfig()

	return nil, nil
}

func initConfig() {

	if resp, err := Client.Get(context.Background(), WEB_CONFIG); err == nil && resp.Count == 0 {
		webConfig, _ := json.Marshal(config.NewInitWebConfig())
		Client.Put(context.Background(), WEB_CONFIG, string(webConfig))
	}
	if resp, err := Client.Get(context.Background(), SERVICE_CONFIG); err == nil && resp.Count == 0 {
		serviceConfig, _ := json.Marshal(config.NewInitServiceConfig())
		Client.Put(context.Background(), SERVICE_CONFIG, string(serviceConfig))
	}
	if resp, err := Client.Get(context.Background(), JUDGE_CONFIG); err == nil && resp.Count == 0 {
		judgeConfig, _ := json.Marshal(config.NewInitJudgeConfig())
		Client.Put(context.Background(), JUDGE_CONFIG, string(judgeConfig))
	}

}
