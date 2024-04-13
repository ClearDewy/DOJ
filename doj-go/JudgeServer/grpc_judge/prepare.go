/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"context"
	"doj-go/JudgeServer/config"
	"doj-go/JudgeServer/internal/redis"
	"doj-go/JudgeServer/internal/sql"
	"github.com/ClearDewy/go-pkg/logrus"
	"os"
	"path/filepath"
	"strconv"
)

func (js *JudgeServer) Prepare() error {

	UseableMutex.Lock()
	Useable = false
	UseableParallelism -= int(js.JudgeItem.Parallelism)
	if UseableParallelism > 0 {
		Useable = true
		go redis.Rdb.RPush(context.Background(), JUDGE_SERVER, config.Conf.JudgeServerAddr)
	}
	UseableMutex.Unlock()
	js.JudgeStatus = &sql.JudgeStatusType{
		Jid:    int(js.JudgeItem.Jid),
		Status: Pending,
	}

	// 获取用户提交的信息
	judgeInfo, err := sql.GetJudgeCode(int(js.JudgeItem.Jid))
	if err != nil {
		logrus.ErrorM(err, "获取提交信息失败")
		return err
	}

	limit, err := sql.GetProblemLanguageLimit(int(js.JudgeItem.Pid), judgeInfo.Lid)
	if err != nil {
		logrus.ErrorM(err, "获取语言时限信息失败")
		return err
	}

	problem, err := sql.GetProblemInfoByPid(int(js.JudgeItem.Pid))
	if err != nil {
		logrus.ErrorM(err, "获取判题题目信息失败")
		return err
	}
	js.JudgeProc = &JudgeProcessType{
		CaseDir:   filepath.Join(config.ROOT_PATH, "problem", strconv.Itoa(problem.Id), string(problem.CaseVersion)),
		UserDir:   filepath.Join(config.ROOT_PATH, "run", strconv.Itoa(js.JudgeStatus.Jid)),
		Problem:   &problem,
		JudgeInfo: &judgeInfo,
		LangLimit: &limit,
		LangCmd:   LanguageList[judgeInfo.Lid],
	}
	err = js.SyncJudgeFile()
	if err != nil {
		logrus.ErrorM(err, "获取判题题目信息失败")
		return err
	}
	return nil
}

func (js *JudgeServer) SyncJudgeFile() error {
	problem := js.JudgeProc.Problem
	dir := filepath.Join(config.ROOT_PATH, "problem", strconv.Itoa(problem.Id), string(problem.CaseVersion))
	_, err := os.Stat(dir)
	// 已经存在
	if err == nil {
		return nil
	}
	// 除不存在以外其他错误
	if !os.IsNotExist(err) {
		return err
	}
	// 不存在，首先删除该题目其他数据版本
	err = os.RemoveAll(filepath.Dir(dir))
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	list, err := sql.GetProblemCaseByPid(problem.Id)
	if err != nil {
		os.RemoveAll(filepath.Dir(dir))
		return err
	}
	for _, value := range list {
		in, err := os.Create(filepath.Join(dir, strconv.Itoa(value.Id)+".in"))
		if err != nil {
			os.RemoveAll(filepath.Dir(dir))
			return err
		}
		_, err = in.Write([]byte(value.Input))
		if err != nil {
			os.RemoveAll(filepath.Dir(dir))
			return err
		}
		in.Close()
		out, err := os.Create(filepath.Join(dir, strconv.Itoa(value.Id)+".out"))
		if err != nil {
			os.RemoveAll(filepath.Dir(dir))
			return err
		}
		_, err = out.Write([]byte(value.Output))
		if err != nil {
			os.RemoveAll(filepath.Dir(dir))
			return err
		}
		out.Close()
	}
	return nil
}
