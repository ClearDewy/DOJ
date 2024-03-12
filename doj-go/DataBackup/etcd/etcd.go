/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package etcd

import (
	"context"
	"doj-go/DataBackup/config"
	"doj-go/DataBackup/utils"
	"encoding/json"
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

func Init(conf *config.Config) {
	var err error
	Client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.EtcdAddr},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    conf.EtcdRootPassword,
	})
	if err != nil {
		utils.HandleError(err, "etcd 连接失败")
	}
	initConfig()
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
