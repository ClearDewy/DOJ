/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package judge

import (
	"context"
	"doj-go/DataBackup/config"
	"doj-go/DataBackup/internal/redis"
	"doj-go/jspb"
	"github.com/ClearDewy/go-pkg/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

type JudgeServerType struct {
	Client jspb.JudgeServerClient
	//Parallelism int
	Error   bool
	Mutex   sync.Mutex
	Conn    *grpc.ClientConn
	LastUse time.Time
	Addr    string
}

var (
	judgeServerConnPool      = make(map[string]*JudgeServerType)
	judgeServerConnPoolMutex = &sync.RWMutex{}
)

func CloseExpireConn() {
	for {
		time.Sleep(config.Conf.JudgeServerExpireTime)
		judgeServerConnPoolMutex.Lock()
		for k, v := range judgeServerConnPool {
			v.Mutex.Lock()
			if v.Conn != nil && time.Now().After(v.LastUse.Add(config.Conf.JudgeServerExpireTime)) {
				err := v.Conn.Close()
				v.Conn = nil
				logrus.ErrorM(err, "CloseExpireConn Error;ip:"+k)
			}
			v.Mutex.Unlock()
		}
		judgeServerConnPoolMutex.Unlock()
	}
}

func ConnectJudgeServer(addr string) (err error) {
	js := &JudgeServerType{
		Addr:    addr,
		LastUse: time.Now(),
	}
	js.Conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.ErrorM(err, "ConnectJudgeServer Error;ip:"+addr)
		return
	}
	js.Client = jspb.NewJudgeServerClient(js.Conn)
	judgeServerConnPoolMutex.Lock()
	judgeServerConnPool[addr] = js
	judgeServerConnPoolMutex.Unlock()
	redis.Rdb.RPush(context.Background(), JUDGE_SERVER, addr)
	return
}

func CheckJudgeServerConn(server *JudgeServerType) {
	server.Mutex.Lock()
	if server.Conn == nil {
		var err error
		server.Conn, err = grpc.Dial(server.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logrus.ErrorM(err, "CheckJudgeServerConn Error;ip:"+server.Addr)
		}

	}
	server.Mutex.Unlock()
}
