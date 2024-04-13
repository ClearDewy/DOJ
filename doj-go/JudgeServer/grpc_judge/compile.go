/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"context"
	"doj-go/JudgeServer/config"
	"doj-go/JudgeServer/sandbox"
	"fmt"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/criyle/go-judge/pb"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (js *JudgeServer) Compile() error {
	// 将提交状态改为编译中
	err := js.UpdateJudgeStatus(Compiling)
	if err != nil {
		return err
	}
	// 编译用户程序
	// 语言需要编译
	if js.JudgeProc.LangCmd.Compile != nil {
		err = js.DefaultCompile()
		if err != nil {
			logrus.ErrorM(err, "编译题目信息失败")
			return err
		}
		if js.JudgeProc.CompileResult.Status != pb.Response_Result_Accepted {
			js.JudgeStatus.Message = string(js.JudgeProc.CompileResult.Files["stderr"])
			return js.UpdateJudgeStatus(CompileError)
		}
		//需要编译即已经编译过了，运行文件为缓存中的结果
		js.JudgeProc.RunCopyInFile = &pb.Request_File{
			File: &pb.Request_File_Cached{
				Cached: &pb.Request_CachedFile{
					FileID: js.JudgeProc.CompileResult.FileIDs[js.JudgeProc.LangCmd.ExePath],
				},
			},
		}
	} else {
		// 如果不需要编译，运行文件就是代码
		js.JudgeProc.RunCopyInFile = &pb.Request_File{
			File: &pb.Request_File_Memory{
				Memory: &pb.Request_MemoryFile{
					Content: []byte(js.JudgeProc.JudgeInfo.Code),
				},
			},
		}
	}

	return nil
}

func (js *JudgeServer) DefaultCompile() error {
	judgeProc := js.JudgeProc
	langCmd := judgeProc.LangCmd
	arg := strings.Replace(langCmd.Compile.Command, "{src_path}", langCmd.SrcPath, 1)
	arg = strings.Replace(arg, "{exe_path}", langCmd.ExePath, 1)

	cmd := &pb.Request_CmdType{
		Args: strings.Split(arg, " "),
		Env:  langCmd.Env,
		Files: []*pb.Request_File{
			{
				File: &pb.Request_File_Memory{
					Memory: &pb.Request_MemoryFile{
						Content: []byte{},
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
		CpuTimeLimit:   uint64(langCmd.MaxCpuTime),
		ClockTimeLimit: uint64(langCmd.MaxRealTime),
		MemoryLimit:    langCmd.MaxMemory,
		StackLimit:     256 << 20,
		ProcLimit:      sandbox.MaxProcLimit,
		CopyIn: map[string]*pb.Request_File{
			langCmd.SrcPath: {
				File: &pb.Request_File_Memory{
					Memory: &pb.Request_MemoryFile{
						Content: []byte(judgeProc.JudgeInfo.Code),
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
				Name:     langCmd.ExePath,
				Optional: true,
			},
		},
	}

	result, err := sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: []*pb.Request_CmdType{cmd},
	})
	if err != nil {
		return err
	}
	judgeProc.CompileResult = result.Results[0]
	return nil
}

// spj或者interactive程序编译
func (js *JudgeServer) SpjOrInteractCompile() error {

	judgeProc := js.JudgeProc
	problem := judgeProc.Problem
	judgeProc.SpjCmd = LanguageList[judgeProc.Problem.SpjLid]
	dir := filepath.Join(config.ROOT_PATH, string(problem.JudgeMode), strconv.Itoa(problem.Id), string(problem.CaseVersion))
	judgeProc.SpjPath = filepath.Join(dir, judgeProc.SpjCmd.ExePath)
	_, err := os.Stat(dir)
	// 已经存在
	if err == nil {
		return nil
	} else {
		// 除不存在以外其他错误
		if !os.IsNotExist(err) {
			return err
		}
		// 不存在，首先删除该题目其他程序版本
		err = os.RemoveAll(filepath.Dir(dir))
		if err != nil {
			return err

		}
		langCmd := LanguageList[problem.SpjLid]
		arg := strings.Replace(langCmd.Compile.Command, "{src_path}", langCmd.SrcPath, 1)
		arg = strings.Replace(arg, "{exe_path}", langCmd.ExePath, 1)

		cmd := &pb.Request_CmdType{
			Args: strings.Split(arg, " "),
			Env:  langCmd.Env,
			Files: []*pb.Request_File{
				{
					File: &pb.Request_File_Memory{
						Memory: &pb.Request_MemoryFile{
							Content: []byte{},
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
			CpuTimeLimit:   uint64(langCmd.MaxCpuTime),
			ClockTimeLimit: uint64(langCmd.MaxRealTime),
			MemoryLimit:    langCmd.MaxMemory,
			StackLimit:     256 << 20,
			ProcLimit:      sandbox.MaxProcLimit,
			CopyIn: map[string]*pb.Request_File{
				langCmd.SrcPath: {
					File: &pb.Request_File_Memory{
						Memory: &pb.Request_MemoryFile{
							Content: []byte(problem.SpjCode),
						},
					},
				},
			},
			CopyOut: []*pb.Request_CmdCopyOutFile{
				{
					Name:     langCmd.ExePath,
					Optional: true,
				},
			},
		}

		result, err := sandbox.Client.Exec(context.Background(), &pb.Request{
			Cmd: []*pb.Request_CmdType{cmd},
		})
		if err != nil {
			return err
		}
		if result.Results[0].Status != pb.Response_Result_Accepted {
			err = fmt.Errorf(result.Results[0].Error)
			logrus.ErrorM(err, "SpjOrInteractCompile Error")
			return err
		}
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			logrus.ErrorM(err, "创建存放"+string(problem.JudgeMode)+"文件夹失败")
			return err
		}
		in, err := os.Create(judgeProc.SpjPath)
		if err != nil {
			logrus.ErrorM(err, "创建"+string(problem.JudgeMode)+"文件失败")
			return err
		}
		defer in.Close()
		_, err = in.Write(result.Results[0].Files[langCmd.ExePath])
		delete(result.Results[0].Files, langCmd.ExePath)

		return err
	}
}
