/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"bufio"
	"context"
	"doj-go/DataBackup/utils"
	"doj-go/JudgeServer/pb"
	"doj-go/JudgeServer/server/grpc_sandbox"
	"doj-go/JudgeServer/server/sql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (js *JudgeServer) Compare(judgeStatus *sql.JudgeStatusType, judgeProc *JudgeProcessType) error {
	// 在之前阶段已经被处理过了
	if judgeStatus.Status != Judging {
		return nil
	}
	for index, id := range judgeProc.Problem.CaseIDs {
		jcs := js.CompareOneDefault(id, judgeProc.RunResults[index])
		err := sql.AddOrUpdateJudgeCaseStatus(&jcs)
		if err != nil {
			return err
		}
	}
	if judgeStatus.Status == Judging {
		judgeStatus.Status = Accepted
	}
	return nil
}

func (js *JudgeServer) CompareOneDefault(pcId int, result *pb.Response_Result) (jcs sql.JudgeCaseStatusType) {
	judgeProc := js.JudgeProc

	jcs.Jid = js.JudgeStatus.Jid
	jcs.ProblemCaseId = pcId
	jcs.Time = result.Time
	jcs.Memory = result.Memory
	jcs.Message = result.Error
	if !checkResultStatus(&jcs, result) {
		return
	}
	// 答案文件
	ans, err := os.Open(filepath.Join(judgeProc.CaseDir, strconv.Itoa(pcId)+".out"))
	if err != nil {
		utils.HandleError(err, "打开答案文件失败")
		jcs.Status = SystemError
		return
	}
	defer ans.Close()
	// 选手答案
	out, err := os.Open(filepath.Join(judgeProc.UserDir, strconv.Itoa(pcId)+".out"))
	if err != nil {
		utils.HandleError(err, "打开选手输出文件失败")
		jcs.Status = SystemError
		return
	}
	defer out.Close()
	ansScan := bufio.NewScanner(ans)
	outScan := bufio.NewScanner(out)
	ansFlag := ansScan.Scan()
	outFlag := outScan.Scan()
	for ansFlag || outFlag {
		ansLine := ""
		outLine := ""

		for ansLine == "" && ansFlag {
			ansLine = ansScan.Text()
			if judgeProc.Problem.IsRemoveEndBlank {
				ansLine = strings.TrimSpace(ansLine)
			}
			ansFlag = ansScan.Scan()
		}
		for outLine == "" && outFlag {
			outLine = outScan.Text()
			if judgeProc.Problem.IsRemoveEndBlank {
				outLine = strings.TrimSpace(outLine)
			}
			outFlag = outScan.Scan()
		}
		if ansLine != outLine {
			jcs.Status = WrongAnswer
			break
		}
	}
	if jcs.Status != WrongAnswer && !ansFlag && !outFlag {
		jcs.Status = Accepted
	}

	checkJudgeCaseStatus(js.JudgeStatus, &jcs)
	return
}
func (js *JudgeServer) CompareOneSpj(pcId int, result *pb.Response_Result) (jcs sql.JudgeCaseStatusType) {
	judgeProc := js.JudgeProc

	jcs.Jid = js.JudgeStatus.Jid
	jcs.ProblemCaseId = pcId
	jcs.Time = result.Time
	jcs.Memory = result.Memory
	jcs.Message = result.Error
	if !checkResultStatus(&jcs, result) {
		return
	}
	// 输入文件
	inPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(pcId)+".in")
	// 答案文件
	ansPath := filepath.Join(judgeProc.CaseDir, strconv.Itoa(pcId)+".out")
	// 选手答案
	outPath := filepath.Join(judgeProc.UserDir, strconv.Itoa(pcId)+".out")
	arg := strings.Replace(judgeProc.SpjCmd.Run.Command, "{exe_path}", judgeProc.SpjCmd.ExePath, 1) + "in out ans"
	cmd := &pb.Request_CmdType{
		Args: strings.Split(arg, " "),
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
		ProcLimit:    grpc_sandbox.MaxProcLimit,
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
	}

	results, err := grpc_sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: []*pb.Request_CmdType{cmd},
	})
	if err != nil {
		utils.HandleError(err, "调用sandbox进行特殊判题失败")
		jcs.Status = SystemError
	} else {
		if results.Results[0].Status == pb.Response_Result_Accepted {
			if results.Results[0].ExitStatus == 0 {
				jcs.Status = Accepted
			} else {
				jcs.Status = WrongAnswer
			}
		} else {
			jcs.Status = SystemError
		}
	}

	checkJudgeCaseStatus(js.JudgeStatus, &jcs)
	return
}

func (js *JudgeServer) CompareOneInteractive(pcId int, userRes *pb.Response_Result, interactRes *pb.Response_Result) (jcs sql.JudgeCaseStatusType) {
	judgeStatus := js.JudgeStatus

	jcs.Jid = judgeStatus.Jid
	jcs.ProblemCaseId = pcId
	jcs.Time = userRes.Time
	jcs.Memory = userRes.Memory
	jcs.Message = userRes.Error
	if !checkResultStatus(&jcs, userRes) || !checkResultStatus(&jcs, interactRes) {
		return
	}

	if userRes.ExitStatus == 0 && interactRes.ExitStatus == 0 {
		judgeStatus.Status = Accepted
	} else {
		judgeStatus.Status = WrongAnswer
	}

	checkJudgeCaseStatus(judgeStatus, &jcs)
	return
}

// 检查sandbox返回结果状态
// 返回为true继续评测，否则结束评测
func checkResultStatus(jcs *sql.JudgeCaseStatusType, result *pb.Response_Result) bool {
	if result.Status == pb.Response_Result_Accepted {
		return true
	}
	switch result.Status {
	case pb.Response_Result_MemoryLimitExceeded:
		jcs.Status = MemoryLimitExceeded
	case pb.Response_Result_TimeLimitExceeded:
		jcs.Status = TimeLimitExceeded
	case pb.Response_Result_InternalError:
		jcs.Status = SystemError
	case pb.Response_Result_OutputLimitExceeded:
		jcs.Status = WrongAnswer

	default:
		jcs.Status = RuntimeError
	}
	return false
}

// 更新总的判题信息
func checkJudgeCaseStatus(judgeStatus *sql.JudgeStatusType, jcs *sql.JudgeCaseStatusType) {
	if jcs.Status == Accepted {
		if judgeStatus.Status == Judging && judgeStatus.Time < jcs.Time {
			judgeStatus.Time = jcs.Time
		}
		if judgeStatus.Status == Judging && judgeStatus.Memory < jcs.Memory {
			judgeStatus.Memory = jcs.Memory
		}
	} else if judgeStatus.Status == Judging {
		judgeStatus.Status = jcs.Status
		judgeStatus.Memory = jcs.Memory
		judgeStatus.Time = jcs.Time
	}
}
