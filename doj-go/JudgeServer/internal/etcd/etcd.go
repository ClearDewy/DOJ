/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package etcd

import (
	"context"
	"doj-go/JudgeServer/config"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"time"
)

var Client *clientv3.Client

func Init() (func() error, func(ctx context.Context) error) {
	etcdConf := clientv3.Config{
		Endpoints:   []string{config.Conf.EtcdAddr},
		DialTimeout: 5 * time.Second,
	}
	if config.Conf.EtcdRootPassword != "" {
		etcdConf.Username = "root"
		etcdConf.Password = config.Conf.EtcdRootPassword
	}
	var err error
	Client, err = clientv3.New(etcdConf)
	if err != nil {
		logrus.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()                          // 在函数结束时调用cancel释放资源
	leaseResp, err := Client.Grant(ctx, 10) // 创建一个10秒的租约
	if err != nil {
		logrus.Fatal(err)
	}
	// cpu核
	_, _ = Client.Put(context.Background(),
		"/judge/cpu-num/"+config.Conf.JudgeServerAddr,
		strconv.Itoa(config.Conf.Parallelism),
		clientv3.WithLease(leaseResp.ID))
	// 同时判题数量核
	_, _ = Client.Put(context.Background(),
		"/judge/parallelism/"+config.Conf.JudgeServerAddr,
		strconv.Itoa(config.Conf.Parallelism),
		clientv3.WithLease(leaseResp.ID))

	keepAliveChan, err := Client.KeepAlive(context.Background(), leaseResp.ID)
	logrus.ErrorM(err, "")

	// 最后更新，避免获取不到前面更新的信息
	// 地址
	_, err = Client.Put(context.Background(),
		"/judge/judge-server/"+config.Conf.JudgeServerAddr,
		config.Conf.JudgeServerName,
		clientv3.WithLease(leaseResp.ID))
	return func() error {
		for {
			select {
			case _, ok := <-keepAliveChan:
				if !ok {
					// 续租失败，可能需要处理
					logrus.ErrorM(err, "")
					return err
				}
				logrus.Info("etcd 续约")
				cpuUsage, _ := cpu.Percent(0, false)
				memUsage, _ := mem.VirtualMemory()
				_, _ = Client.Put(context.Background(),
					"/judge/cpu-usage/"+config.Conf.JudgeServerAddr,
					strconv.FormatFloat(cpuUsage[0], 'f', 1, 32),
					clientv3.WithLease(leaseResp.ID))
				_, _ = Client.Put(context.Background(),
					"/judge/mem-usage/"+config.Conf.JudgeServerAddr,
					strconv.FormatFloat(memUsage.UsedPercent, 'f', 1, 32),
					clientv3.WithLease(leaseResp.ID))
			}
		}
	}, stop
}

func stop(ctx context.Context) error {
	Client.Delete(context.Background(), "/judge/judge-server/"+config.Conf.JudgeServerAddr)
	Client.Delete(context.Background(), "/judge/cpu-num/"+config.Conf.JudgeServerAddr)
	Client.Delete(context.Background(), "/judge/parallelism/"+config.Conf.JudgeServerAddr)
	Client.Delete(context.Background(), "/judge/cpu-usage/"+config.Conf.JudgeServerAddr)
	Client.Delete(context.Background(), "/judge/mem-usage/"+config.Conf.JudgeServerAddr)
	return nil
}
