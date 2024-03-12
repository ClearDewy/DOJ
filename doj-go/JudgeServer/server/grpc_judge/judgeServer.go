/**
 * @ Author: ClearDewy
 * @ Desc: 实现gRPC代码
 **/
package grpc_judge

import (
	"context"
	"doj-go/JudgeServer/pb"
	"doj-go/JudgeServer/server/config"
	"doj-go/JudgeServer/server/grpc_sandbox"
	"doj-go/JudgeServer/server/sql"
	"doj-go/JudgeServer/server/utils"
	"doj-go/jpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"path/filepath"
	"strconv"
)

type JudgeServer struct {
	jpb.UnimplementedJudgeServerServer
	JudgeStatus *sql.JudgeStatusType
	JudgeProc   *JudgeProcessType
}

func test() {
	cmd := &pb.Request_CmdType{
		Args: []string{"/usr/bin/g++", "a.cc", "-o", "a"},
		Env:  []string{"PATH=/usr/bin:/bin"},
		Files: []*pb.Request_File{
			{
				File: &pb.Request_File_Local{
					Local: &pb.Request_LocalFile{
						Src: "/doj/judge/problem/1000/0/1.in",
					},
				},
			},
			{
				File: &pb.Request_File_Pipe{
					&pb.Request_PipeCollector{
						Name: "stdout",
						Max:  10240,
					},
				},
			},
			{
				File: &pb.Request_File_Pipe{
					&pb.Request_PipeCollector{
						Name: "stderr",
						Max:  10240,
					},
				},
			},
		},
		CpuTimeLimit: 10000000000,
		MemoryLimit:  104857600,
		ProcLimit:    128,
		CopyIn: map[string]*pb.Request_File{
			"a.cc": {
				File: &pb.Request_File_Memory{
					Memory: &pb.Request_MemoryFile{
						Content: []byte("#include <iostream>\nusing namespace std;\nint main() {\nint a, b;\ncin >> a >> b;\ncout << a + b << endl;\n}"),
					},
				},
			},
		},
		CopyOut: []*pb.Request_CmdCopyOutFile{
			{
				Name: "stdout",
			},
			{
				Name: "stderr",
			},
		},
		CopyOutCached: []*pb.Request_CmdCopyOutFile{
			{
				Name:     "a",
				Optional: true,
			},
		},
	}
	result, err := grpc_sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: []*pb.Request_CmdType{cmd},
	})
	if err != nil {
		utils.HandleError(err, "")
		return
	}
	defer grpc_sandbox.DeleteFileByResult(result.Results[0])
	logrus.Info(result)
}

func (js *JudgeServer) Judge(c context.Context, judge_item *jpb.JudgeItem) (*emptypb.Empty, error) {
	js.JudgeStatus = &sql.JudgeStatusType{
		Jid:    int(judge_item.Jid),
		Status: Pending,
	}
	problem, err := sql.GetProblemInfoByPid(int(judge_item.Pid))

	if err != nil {
		utils.HandleError(err, "获取判题题目信息失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	js.JudgeProc = &JudgeProcessType{
		CaseDir: filepath.Join(config.ROOT_PATH, "problem", strconv.Itoa(problem.Id), string(problem.CaseVersion)),
		UserDir: filepath.Join(config.ROOT_PATH, "run", strconv.Itoa(js.JudgeStatus.Jid)),
		Problem: &problem,
	}
	err = js.SyncJudgeFile()
	if err != nil {
		utils.HandleError(err, "获取判题题目信息失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}

	// 编译spj程序或者交互程序
	if js.JudgeProc.Problem.JudgeMode != DEFAULT {
		result, err := js.SpjOrInteractCompile()
		if err != nil {
			utils.HandleError(err, "编译题目信息失败")
		}
		if err != nil || result.Status != pb.Response_Result_Accepted {
			js.JudgeStatus.Status = SystemError
			js.UpdateJudgeStatus()
			return nil, err
		}
	}
	if err != nil {
		utils.HandleError(err, "编译题目特殊程序失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	// 将提交状态改为编译中
	js.JudgeStatus.Status = Compiling
	err = js.UpdateJudgeStatus()
	if err != nil {
		utils.HandleError(err, "更新判题信息失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	// 获取用户提交的信息
	judgeInfo, err := sql.GetJudgeCode(int(judge_item.Jid))
	js.JudgeProc.JudgeInfo = &judgeInfo
	if err != nil {
		utils.HandleError(err, "获取提交信息失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	// 编译用户程序
	js.JudgeProc.LangCmd = LanguageList[js.JudgeProc.JudgeInfo.Lid]
	// 语言需要编译
	if js.JudgeProc.LangCmd.Compile != nil {
		err := js.DefaultCompile()
		defer grpc_sandbox.DeleteFileByResult(js.JudgeProc.CompileResult)
		if err != nil {
			utils.HandleError(err, "编译题目信息失败")
			js.JudgeStatus.Status = SystemError
			js.UpdateJudgeStatus()
			return nil, err
		}
		if js.JudgeProc.CompileResult.Status != pb.Response_Result_Accepted {
			js.JudgeStatus.Status = CompileError
			js.JudgeStatus.Message = string(js.JudgeProc.CompileResult.Files["stderr"])
			js.UpdateJudgeStatus()
			return nil, nil
		}
	}

	// 将提交状态改为运行中
	js.JudgeStatus.Status = Judging
	err = js.UpdateJudgeStatus()
	if err != nil {
		utils.HandleError(err, "更新判题信息失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	limit, err := sql.GetProblemLanguageLimit(int(judge_item.Pid), judgeInfo.Lid)
	if err != nil {
		utils.HandleError(err, "获取语言时限信息失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	js.JudgeProc.LangLimit = &limit
	err = js.Run()
	if err != nil {
		utils.HandleError(err, "运行程序时系统出现错误")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	// 删除运行产生的临时文件
	defer os.RemoveAll(js.JudgeProc.UserDir)

	err = js.Compare(js.JudgeStatus, js.JudgeProc)
	if err != nil {
		utils.HandleError(err, "添加题目数据评测结果失败")
		js.JudgeStatus.Status = SystemError
		js.UpdateJudgeStatus()
		return nil, err
	}
	js.UpdateJudgeStatus()
	return nil, nil
}

func (js *JudgeServer) UpdateJudgeStatus() error {
	return sql.UpdateJudgeStatus(js.JudgeStatus)
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
