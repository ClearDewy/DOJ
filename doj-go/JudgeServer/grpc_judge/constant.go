/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	sql2 "doj-go/JudgeServer/internal/sql"
	"github.com/criyle/go-judge/pb"
)

const (
	Judging                = -4
	Compiling              = -3
	Pending                = -2
	Submitting             = -1
	Accepted               = 1
	WrongAnswer            = 2
	TimeLimitExceeded      = 3
	MemoryLimitExceeded    = 4
	CompileError           = 5
	RuntimeError           = 6
	PartialAccepted        = 7
	SystemError            = 8
	PresentationError      = 9
	SubmittedFailed        = 10
	NotSubmitted           = 11
	Cancelled              = 12
	SubmittedUnknownResult = 13
)

const (
	// 题目评测模式
	DEFAULT     = "default"
	SPJ         = "spj"
	INTERACTIVE = "interactive"

	// 样例评测模式
	//DEFAULT     = "default"
	ERGODIC_WITHOUT_ERROR = "ergodic_without_error"
	SUBTASK_LOWEST        = "subtask_lowest"
	SUBTASK_AVERAGE       = "subtask_average"
)

type JudgeProcessType struct {
	LangCmd       *LanguageType
	Problem       *sql2.ProblemJudge
	JudgeInfo     *sql2.JudgeInfoType
	LangLimit     *sql2.LanguageLimitType
	CaseDir       string
	UserDir       string
	RunCopyInFile *pb.Request_File
	CompileResult *pb.Response_Result
	RunResults    []*pb.Response_Result
	SpjCmd        *LanguageType
	SpjPath       string
}
