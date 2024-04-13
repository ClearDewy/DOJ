/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"context"
	"doj-go/JudgeServer/config"
	"doj-go/JudgeServer/sandbox"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/criyle/go-judge/pb"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (js *JudgeServer) Run() error {
	// 将提交状态改为运行中
	err := js.UpdateJudgeStatus(Judging)
	if err != nil {
		return err
	}

	judgeProc := js.JudgeProc
	judgeStatus := js.JudgeStatus

	// 出了遇错即止，其他为评测全部
	// 但是需要控制同时最大评测数量
	if judgeProc.Problem.JudgeCaseMode == ERGODIC_WITHOUT_ERROR {

	}
	// 创建用户允许文件输出文件夹
	err = os.MkdirAll(filepath.Join(config.ROOT_PATH, "run", strconv.Itoa(judgeStatus.Jid)), 0777)
	switch judgeProc.Problem.JudgeMode {
	case DEFAULT:
		return js.RunAllDefault()
	case INTERACTIVE:
		return js.RunAllInteract()
	default:
		logrus.Error("未知的评测模式")
	}

	return nil
}

func (js *JudgeServer) RunAllDefault() error {
	judgeProc := js.JudgeProc
	countCase := len(judgeProc.Problem.CaseIDs)
	runCmds := make([]*pb.Request_CmdType, 0, countCase)
	checkCmds := make([]*pb.Request_CmdType, 0, countCase)
	for i := 0; i < countCase; i++ {
		// 用户程序的运行命令
		userArg := strings.Replace(judgeProc.LangCmd.Run.Command, "{exe_path}", judgeProc.LangCmd.ExePath, 1)
		// 判题程序
		checkArg := strings.Replace(judgeProc.SpjCmd.Run.Command, "{exe_path}", judgeProc.SpjCmd.ExePath, 1) + " in out ans"
		// 输入文件
		inPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".in")
		// 答案文件
		ansPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")
		// 选手答案
		outPath := filepath.Join(judgeProc.UserDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")

		runCmds = append(runCmds, &pb.Request_CmdType{
			Args:         strings.Split(userArg, " "),
			Env:          judgeProc.LangCmd.Env,
			CpuTimeLimit: judgeProc.LangLimit.TimeLimit,
			MemoryLimit:  judgeProc.LangLimit.MemoryLimit,
			StackLimit:   judgeProc.LangLimit.StackLimit,
			Files: []*pb.Request_File{
				{
					File: &pb.Request_File_Local{
						Local: &pb.Request_LocalFile{
							Src: inPath,
						},
					},
				},
				{
					File: &pb.Request_File_Pipe{
						Pipe: &pb.Request_PipeCollector{
							Name: "stdout",
							Max:  int64(config.Conf.CopyOutLimit.Byte()),
						},
					},
				},
				{
					File: &pb.Request_File_Pipe{
						Pipe: &pb.Request_PipeCollector{
							Name: "stderr",
							Max:  10240,
						},
					},
				},
			},
			CopyIn: map[string]*pb.Request_File{
				judgeProc.LangCmd.ExePath: judgeProc.RunCopyInFile,
			},
			CopyOut: []*pb.Request_CmdCopyOutFile{
				{
					Name:     "stderr",
					Optional: true,
				},
			},
		})

		checkCmds = append(checkCmds, &pb.Request_CmdType{
			Args: strings.Split(checkArg, " "),
			Env:  judgeProc.SpjCmd.Env,
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
			CpuTimeLimit: uint64(judgeProc.SpjCmd.MaxCpuTime),
			MemoryLimit:  judgeProc.SpjCmd.MaxMemory,
			StackLimit:   256 << 20,
			ProcLimit:    sandbox.MaxProcLimit,
			CopyIn: map[string]*pb.Request_File{
				judgeProc.SpjCmd.ExePath: {
					File: &pb.Request_File_Local{
						Local: &pb.Request_LocalFile{
							Src: judgeProc.SpjPath,
						},
					},
				},
				"in": {
					File: &pb.Request_File_Local{
						Local: &pb.Request_LocalFile{
							Src: inPath,
						},
					},
				},
				"out": {
					File: &pb.Request_File_Local{
						Local: &pb.Request_LocalFile{
							Src: outPath,
						},
					},
				},
				"ans": {
					File: &pb.Request_File_Local{
						Local: &pb.Request_LocalFile{
							Src: ansPath,
						},
					},
				},
			},
			CopyOut: []*pb.Request_CmdCopyOutFile{
				{
					Name:     "stderr",
					Optional: true,
				},
			},
		})
	}
	results, err := sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: runCmds,
	})
	if err != nil {
		logrus.ErrorM(err, "运行程序时发生错误")
		return err
	}
	js.JudgeProc.RunResults = results.Results

	for i := 0; i < countCase; i++ {
		// 选手答案
		outPath := filepath.Join(judgeProc.UserDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")
		file, err := os.Create(outPath)
		if err != nil {
			logrus.ErrorM(err, "创建用户输出文件失败")
			return err
		}
		_, err = file.Write(results.Results[i].Files["stdout"])
		delete(results.Results[i].Files, "stdout")
		file.Close()
		if err != nil {
			return err
		}
	}

	results, err = sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: checkCmds,
	})
	if err != nil {
		logrus.ErrorM(err, "运行检查程序时发生错误")
		return err
	}
	js.JudgeProc.CheckResults = results.Results
	return nil
}

func (js *JudgeServer) RunAllInteract() error {
	judgeProc := js.JudgeProc
	countCase := len(judgeProc.Problem.CaseIDs)
	cmds := make([]*pb.Request_CmdType, 0, countCase<<1)
	pips := make([]*pb.Request_PipeMap, 0, countCase<<1)
	for i := 0; i < countCase; i++ {
		// 用户程序的运行命令
		userArg := strings.Replace(judgeProc.LangCmd.Run.Command, "{exe_path}", judgeProc.LangCmd.ExePath, 1)
		spjArg := strings.Replace(judgeProc.SpjCmd.Run.Command, "{exe_path}", judgeProc.SpjCmd.ExePath, 1) + " in out ans"
		// 输入文件
		inPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".in")
		// 答案文件
		ansPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")
		// 选手答案
		outPath := filepath.Join(judgeProc.UserDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")

		// 创建选手输出的空文件，不创建会报错，选手不会输出到文件中
		file, err := os.Create(outPath)
		if err != nil {
			return err
		}
		file.Close()

		cmds = append(cmds, &pb.Request_CmdType{
			Args:         strings.Split(userArg, " "),
			Env:          judgeProc.LangCmd.Env,
			CpuTimeLimit: judgeProc.LangLimit.TimeLimit,
			MemoryLimit:  judgeProc.LangLimit.MemoryLimit,
			StackLimit:   judgeProc.LangLimit.StackLimit,
			Files: []*pb.Request_File{
				nil, nil,
				{
					File: &pb.Request_File_Pipe{
						&pb.Request_PipeCollector{
							Name: "stderr",
							Max:  10240,
						},
					},
				},
			},
			CopyIn: map[string]*pb.Request_File{
				judgeProc.LangCmd.ExePath: judgeProc.RunCopyInFile,
			},
			CopyOut: []*pb.Request_CmdCopyOutFile{
				{
					Name:     "stderr",
					Optional: true,
				},
			},
		},
			&pb.Request_CmdType{
				Args: strings.Split(spjArg, " "),
				Env:  judgeProc.SpjCmd.Env,
				Files: []*pb.Request_File{
					nil, nil,
					{
						File: &pb.Request_File_Pipe{
							&pb.Request_PipeCollector{
								Name: "stderr",
								Max:  10240,
							},
						},
					},
				},
				CpuTimeLimit: uint64(judgeProc.SpjCmd.MaxCpuTime),
				MemoryLimit:  judgeProc.SpjCmd.MaxMemory,
				StackLimit:   256 << 20,
				CopyIn: map[string]*pb.Request_File{
					judgeProc.SpjCmd.ExePath: {
						File: &pb.Request_File_Local{
							Local: &pb.Request_LocalFile{
								Src: judgeProc.SpjPath,
							},
						},
					},
					"in": {
						File: &pb.Request_File_Local{
							Local: &pb.Request_LocalFile{
								Src: inPath,
							},
						},
					},
					"out": {
						File: &pb.Request_File_Local{
							Local: &pb.Request_LocalFile{
								Src: outPath,
							},
						},
					},
					"ans": {
						File: &pb.Request_File_Local{
							Local: &pb.Request_LocalFile{
								Src: ansPath,
							},
						},
					},
				},
				CopyOut: []*pb.Request_CmdCopyOutFile{
					{
						Name:     "stderr",
						Optional: true,
					},
				},
			})
		pips = append(pips, &pb.Request_PipeMap{
			// 选手程序输出端
			In: &pb.Request_PipeMap_PipeIndex{
				Index: int32(i << 1),
				Fd:    1,
			},
			// interact 输入端
			Out: &pb.Request_PipeMap_PipeIndex{
				Index: int32((i << 1) + 1),
				Fd:    0,
			},
		},
			&pb.Request_PipeMap{
				// interact 输出端
				In: &pb.Request_PipeMap_PipeIndex{
					Index: int32((i << 1) + 1),
					Fd:    1,
				},
				// 选手程序输入端
				Out: &pb.Request_PipeMap_PipeIndex{
					Index: int32(i << 1),
					Fd:    0,
				},
			})
	}

	results, err := sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd:         cmds,
		PipeMapping: pips,
	})
	if err != nil {
		return err
	}

	for i := 0; i < countCase; i++ {
		judgeProc.RunResults = append(judgeProc.RunResults, results.Results[i<<1])
		judgeProc.CheckResults = append(judgeProc.CheckResults, results.Results[(i<<1)+1])
	}

	return nil
}
