/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"context"
	"doj-go/DataBackup/utils"
	"doj-go/JudgeServer/pb"
	"doj-go/JudgeServer/server/config"
	"doj-go/JudgeServer/server/grpc_sandbox"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
		ProcLimit:      grpc_sandbox.MaxProcLimit,
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

	result, err := grpc_sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: []*pb.Request_CmdType{cmd},
	})
	if err != nil {
		return err
	}
	judgeProc.CompileResult = result.Results[0]
	return nil
}

// spj或者interactive程序编译
func (js *JudgeServer) SpjOrInteractCompile() (*pb.Response_Result, error) {
	judgeProc := js.JudgeProc
	problem := judgeProc.Problem
	judgeProc.SpjCmd = LanguageList[judgeProc.Problem.SpjLid]
	dir := filepath.Join(config.ROOT_PATH, string(problem.JudgeMode), strconv.Itoa(problem.Id), string(problem.CaseVersion))
	judgeProc.SpjPath = filepath.Join(dir, judgeProc.SpjCmd.ExePath)
	_, err := os.Stat(dir)
	// 已经存在
	if err == nil {
		return &pb.Response_Result{
			Status: pb.Response_Result_Accepted,
		}, err
	} else {
		// 除不存在以外其他错误
		if !os.IsNotExist(err) {
			return nil, err
		}
		// 不存在，首先删除该题目其他程序版本
		err = os.RemoveAll(filepath.Dir(dir))
		if err != nil {
			return nil, err

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
			ProcLimit:      grpc_sandbox.MaxProcLimit,
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

		result, err := grpc_sandbox.Client.Exec(context.Background(), &pb.Request{
			Cmd: []*pb.Request_CmdType{cmd},
		})
		if err != nil {
			return nil, err
		}
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			utils.HandleError(err, "创建存放"+string(problem.JudgeMode)+"文件夹失败")
			return result.Results[0], nil
		}
		in, err := os.Create(judgeProc.SpjPath)
		if err != nil {
			utils.HandleError(err, "创建"+string(problem.JudgeMode)+"文件失败")
			return result.Results[0], nil
		}
		defer in.Close()
		in.Write(result.Results[0].Files[string(problem.JudgeMode)])

		return result.Results[0], nil
	}
}
