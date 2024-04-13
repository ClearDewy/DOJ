/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"doj-go/JudgeServer/internal/sql"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/criyle/go-judge/pb"
)

func (js *JudgeServer) Check() error {
	// 在之前阶段已经被处理过了
	judgeProc := js.JudgeProc
	err := js.UpdateJudgeStatus(Judging)
	if err != nil {
		return err
	}

	for i := 0; i < len(judgeProc.RunResults); i++ {
		runResult := judgeProc.RunResults[i]
		checkResult := judgeProc.CheckResults[i]
		jcs := &sql.JudgeCaseStatusType{
			Status:        Judging,
			Jid:           js.JudgeStatus.Jid,
			ProblemCaseId: judgeProc.Problem.CaseIDs[i],
			Time:          runResult.Time,
			Memory:        runResult.Memory,
			Message:       runResult.Error,
		}
		if updateJudgeCaseStatus(jcs, runResult) && updateJudgeCaseStatus(jcs, checkResult) {
			jcs.Status = Accepted
		}

		updateJudgeStatus(js.JudgeStatus, jcs)
		err := sql.AddOrUpdateJudgeCaseStatus(jcs)
		if err != nil {
			logrus.ErrorM(err, "更新测试数据出错")
			return err
		}
	}

	// 结果判断完了没有出错
	if js.JudgeStatus.Status == Judging {
		js.JudgeStatus.Status = Accepted
	}

	return sql.UpdateJudgeStatus(js.JudgeStatus)
}

// 检查sandbox返回结果状态
// 返回为true继续评测，否则结束评测
func updateJudgeCaseStatus(jcs *sql.JudgeCaseStatusType, result *pb.Response_Result) bool {
	if result.Status == pb.Response_Result_Accepted {
		return true
	}
	switch result.Status {
	case pb.Response_Result_MemoryLimitExceeded:
		jcs.Status = MemoryLimitExceeded
	case pb.Response_Result_TimeLimitExceeded:
		jcs.Status = TimeLimitExceeded
	case pb.Response_Result_OutputLimitExceeded:
		jcs.Status = WrongAnswer

	default:
		jcs.Status = RuntimeError
	}
	return false
}

// 更新总的判题信息
func updateJudgeStatus(judgeStatus *sql.JudgeStatusType, jcs *sql.JudgeCaseStatusType) {
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
