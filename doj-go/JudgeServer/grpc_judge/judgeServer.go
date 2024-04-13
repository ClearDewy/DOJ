/**
 * @ Author: ClearDewy
 * @ Desc: 实现gRPC代码
 **/
package grpc_judge

import (
	"context"
	"doj-go/JudgeServer/config"
	"doj-go/JudgeServer/internal/redis"
	"doj-go/JudgeServer/internal/sql"
	"doj-go/JudgeServer/sandbox"
	"doj-go/jspb"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"sync"
)

type JudgeServer struct {
	jspb.UnimplementedJudgeServerServer
	JudgeStatus *sql.JudgeStatusType
	JudgeProc   *JudgeProcessType
	JudgeItem   *jspb.JudgeItem
}

var (
	Useable            = true
	UseableMutex       = &sync.RWMutex{}
	UseableParallelism = config.Conf.Parallelism
)

func (js *JudgeServer) Judge(c context.Context, judge_item *jspb.JudgeItem) (*emptypb.Empty, error) {
	js.JudgeItem = judge_item

	for _, step := range []func() error{
		js.Prepare,
		js.Compile,
		js.SpjOrInteractCompile,
		js.Run,
		js.Check,
		js.Recovery,
	} {
		err := step()
		if err != nil {
			js.UpdateJudgeStatus(SystemError)
			break
		}
	}
	return nil, nil
}

func (js *JudgeServer) UpdateJudgeStatus(status int) error {
	js.JudgeStatus.Status = status
	return sql.UpdateJudgeStatus(js.JudgeStatus)
}

func (js *JudgeServer) Recovery() error {
	// 编译过
	if js.JudgeProc.LangCmd.Compile != nil {
		sandbox.DeleteFileByResult(js.JudgeProc.CompileResult)
	}
	// sandbox.DeleteFileByResult(js.JudgeProc.CheckCompileResult)
	UseableMutex.Lock()
	UseableParallelism += int(js.JudgeItem.Parallelism)
	if UseableParallelism > 0 && !Useable {
		Useable = true
		go redis.Rdb.RPush(context.Background(), JUDGE_SERVER, config.Conf.JudgeServerAddr)
	}
	UseableMutex.Unlock()

	// 删除运行产生的临时文件
	defer os.RemoveAll(js.JudgeProc.UserDir)
	return nil
}
