/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sql

import (
	"database/sql"
	dSql "github.com/ClearDewy/go-pkg/sql"
)

type ProblemJudge struct {
	Id               int
	Type             int
	Oj               dSql.String
	CaseVersion      dSql.String
	SpjLid           dSql.Int
	SpjCode          dSql.String
	JudgeMode        dSql.String
	JudgeCaseMode    dSql.String
	IsRemoveEndBlank dSql.Bool
	UserExtraFile    dSql.String
	JudgeExtraFile   dSql.String
	IsFileIo         dSql.Bool
	IoReadFileName   dSql.String
	IoWriteFileName  dSql.String
	CaseIDs          []int
}

type ProblemCase struct {
	Id     int
	Input  dSql.String
	Output dSql.String
}

func GetProblemInfoByPid(pid int) (problem ProblemJudge, err error) {
	row := db.QueryRow("SELECT `id`,`type`,`oj`,`case_version`,`spj_lid`,`spj_code`,`judge_mode`,`judge_case_mode`,`is_remove_end_blank`,`user_extra_file`,`judge_extra_file`,`is_file_io`,`io_write_file_name`,`io_read_file_name` FROM problem WHERE id=?", pid)
	err = row.Scan(&problem.Id, &problem.Type, &problem.Oj, &problem.CaseVersion, &problem.SpjLid, &problem.SpjCode,
		&problem.JudgeMode, &problem.JudgeCaseMode, &problem.IsRemoveEndBlank, &problem.UserExtraFile, &problem.JudgeExtraFile, &problem.IsFileIo, &problem.IoReadFileName, &problem.IoWriteFileName)
	if err != nil {
		return
	}
	problem.CaseIDs, err = GetProblemCaseIDs(pid)
	return
}

func GetProblemCaseByPid(pid int) (problemCaseList []*ProblemCase, err error) {
	var rows *sql.Rows
	rows, err = db.Query("SELECT `id`,`input`,`output` FROM problem_case WHERE pid=?", pid)

	for rows.Next() {
		problemCase := ProblemCase{}
		err = rows.Scan(&problemCase.Id, &problemCase.Input, &problemCase.Output)
		if err != nil {
			return
		}
		problemCaseList = append(problemCaseList, &problemCase)
	}

	return
}

func GetProblemCaseIDs(pid int) ([]int, error) {
	caseIDs := make([]int, 0)
	rows, err := db.Query("SELECT `id` FROM problem_case WHERE pid=?", pid)
	if err != nil {
		return nil, err
	}
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		caseIDs = append(caseIDs, id)
	}
	return caseIDs, nil
}
