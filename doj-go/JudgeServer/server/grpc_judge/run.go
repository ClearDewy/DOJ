/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"context"
	"doj-go/DataBackup/judge"
	"doj-go/DataBackup/utils"
	"doj-go/JudgeServer/pb"
	"doj-go/JudgeServer/server/config"
	"doj-go/JudgeServer/server/grpc_sandbox"
	"doj-go/JudgeServer/server/sql"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func (js *JudgeServer) Run() error {
	judgeProc := js.JudgeProc
	judgeStatus := js.JudgeStatus
	if judgeProc.LangCmd.Compile == nil {
		// 如果不需要编译，运行文件就是代码
		judgeProc.RunCopyInFile = &pb.Request_File{
			File: &pb.Request_File_Memory{
				Memory: &pb.Request_MemoryFile{
					Content: []byte(judgeProc.JudgeInfo.Code),
				},
			},
		}
	} else {
		//需要编译即已经编译过了，运行文件为缓存中的结果
		judgeProc.RunCopyInFile = &pb.Request_File{
			File: &pb.Request_File_Cached{
				Cached: &pb.Request_CachedFile{
					FileID: judgeProc.CompileResult.FileIDs[judgeProc.LangCmd.ExePath],
				},
			},
		}
	}

	part_num := runtime.NumCPU()
	countCase := len(judgeProc.Problem.CaseIDs)

	// 出了遇错即止，其他为评测全部
	// 但是需要控制同时最大评测数量
	if judgeProc.Problem.JudgeCaseMode == ERGODIC_WITHOUT_ERROR {
		part_num = 1
	}

	judgeProc.RunResults = make([]*pb.Response_Result, 0, countCase)
	// 创建用户允许文件输出文件夹
	err := os.MkdirAll(filepath.Join(config.ROOT_PATH, "run", strconv.Itoa(judgeStatus.Jid)), 0777)
	if err != nil {
		return err
	}
	for i := 0; i < countCase; i += part_num {
		r := i + part_num
		if r > countCase {
			r = countCase
		}
		var partResults []*pb.Response_Result

		// 如果是交互判题
		if judgeProc.Problem.JudgeMode == INTERACTIVE {
			partResults, err = js.RunPartInteract(i, r)
		} else {
			partResults, err = js.RunPart(i, r)
		}
		if err != nil {
			return err
		}
		judgeProc.RunResults = append(judgeProc.RunResults, partResults...)
		// 如果是遇错即止就直接判断结果
		if part_num == 1 {
			var jcs sql.JudgeCaseStatusType
			switch judgeProc.Problem.JudgeMode {
			case DEFAULT:
				jcs = js.CompareOneDefault(judgeProc.Problem.CaseIDs[i], partResults[0])
			case SPJ:
				jcs = js.CompareOneSpj(judgeProc.Problem.CaseIDs[i], partResults[0])
			case INTERACTIVE:
				jcs = js.CompareOneInteractive(judgeProc.Problem.CaseIDs[i], partResults[0], partResults[1])
			default:
				logrus.Error("未知的判题模式")
			}
			err = sql.AddOrUpdateJudgeCaseStatus(&jcs)
			if err != nil {
				utils.HandleError(err, "添加题目数据评测结果失败")
				judgeStatus.Status = judge.SystemError
				sql.UpdateJudgeStatus(judgeStatus)
				return err
			}
		}
	}
	// 如果是遇错即止就直接判断结果
	if part_num == 1 && judgeStatus.Status == judge.Judging {
		judgeStatus.Status = judge.Accepted
	}
	return nil
}

func (js *JudgeServer) RunPart(l, r int) ([]*pb.Response_Result, error) {
	judgeProc := js.JudgeProc
	cmds := make([]*pb.Request_CmdType, 0, r-l)
	for i := l; i < r; i++ {
		inPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".in")
		ansPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")
		// 查看std输出的大小
		stat, err := os.Stat(ansPath)
		if err != nil {
			return nil, err
		}
		maxOutPutLimit := config.Conf.OutputLimit
		if stat.Size() > maxOutPutLimit {
			maxOutPutLimit = stat.Size()
		}
		arg := strings.Replace(judgeProc.LangCmd.Run.Command, "{exe_path}", judgeProc.LangCmd.ExePath, 1)
		cmds = append(cmds, &pb.Request_CmdType{
			Args:         strings.Split(arg, " "),
			Env:          judgeProc.LangCmd.Env,
			CpuTimeLimit: judgeProc.LangLimit.TimeLimit,
			MemoryLimit:  judgeProc.LangLimit.MemoryLimit,
			StackLimit:   judgeProc.LangLimit.StackLimit,
			Files: []*pb.Request_File{
				{
					File: &pb.Request_File_Local{
						&pb.Request_LocalFile{
							Src: inPath,
						},
					},
				},
				{
					File: &pb.Request_File_Pipe{
						&pb.Request_PipeCollector{
							Name: "stdout",
							// 最大输出设置为 ans 文件的两倍
							Max: maxOutPutLimit,
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
			CopyIn: map[string]*pb.Request_File{
				judgeProc.LangCmd.ExePath: judgeProc.RunCopyInFile,
			},
			CopyOut: []*pb.Request_CmdCopyOutFile{
				{
					Name:     "stdout",
					Optional: true,
				},
				{
					Name:     "stderr",
					Optional: true,
				},
			},
		})
	}
	results, err := grpc_sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: cmds,
	})
	if err != nil {
		return nil, err
	}
	for index, result := range results.Results {
		outPath := filepath.Join(judgeProc.UserDir, strconv.Itoa(judgeProc.Problem.CaseIDs[index+l])+".out")
		out, err := os.Create(outPath)
		if err != nil {
			return nil, err
		}
		_, err = out.Write(result.Files["stdout"])
		delete(result.Files, "stdout")
		if err != nil {
			return nil, err
		}
		out.Close()
	}
	return results.Results, nil
}

func (js *JudgeServer) RunPartInteract(l, r int) ([]*pb.Response_Result, error) {
	judgeProc := js.JudgeProc
	cmds := make([]*pb.Request_CmdType, 0, (r-l)<<1)
	pips := make([]*pb.Request_PipeMap, 0, (r-l)<<1)
	for i := l; i < r; i++ {
		// 用户程序的运行命令
		userArg := strings.Replace(judgeProc.LangCmd.Run.Command, "{exe_path}", judgeProc.LangCmd.ExePath, 1)
		spjArg := strings.Replace(judgeProc.SpjCmd.Run.Command, "{exe_path}", judgeProc.SpjCmd.ExePath, 1) + "in out ans"
		// 输入文件
		inPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".in")
		// 答案文件
		ansPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")
		// 选手答案
		outPath := filepath.Join(judgeProc.UserDir, strconv.Itoa(judgeProc.Problem.CaseIDs[i])+".out")

		// 创建选手输出的空文件，不创建会报错，选手不会输出到文件中
		file, err := os.Create(outPath)
		if err != nil {
			return nil, err
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
				Index: int32((i - l) * 2),
				Fd:    1,
			},
			// interact 输入端
			Out: &pb.Request_PipeMap_PipeIndex{
				Index: int32((i-l)*2 + 1),
				Fd:    0,
			},
		},
			&pb.Request_PipeMap{
				// interact 输出端
				In: &pb.Request_PipeMap_PipeIndex{
					Index: int32((i-l)*2 + 1),
					Fd:    1,
				},
				// 选手程序输入端
				Out: &pb.Request_PipeMap_PipeIndex{
					Index: int32((i - l) * 2),
					Fd:    0,
				},
			})
	}

	results, err := grpc_sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd:         cmds,
		PipeMapping: pips,
	})
	if err != nil {
		return nil, err
	}
	return results.Results, nil
}
