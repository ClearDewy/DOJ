/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sql

import "doj-go/DataBackup/judge"

type AddJudgeType struct {
	Uid       string `json:"uid"`
	Username  string `json:"username"`
	Pid       int    `json:"pid"`
	ProblemId string `json:"problem_id"`
	Code      string `json:"code"`
	Lid       int    `json:"lid"`
}
type ProblemJudgeType struct {
	Pid         int `json:"pid"`
	Parallelism int `json:"countCase"`
}

func AddJudge(spj *AddJudgeType) (int64, error) {
	resp, err := db.Exec("INSERT INTO judge(`uid`,`username`,`pid`,`problem_id`,`code`,`lid`,`submit_time`,`status`,`length`) VALUES (?,?,?,?,?,?,NOW(),?,?)",
		spj.Uid, spj.Username, spj.Pid, spj.ProblemId, spj.Code, spj.Lid, judge.Pending, len(spj.Code))
	if err != nil {
		return 0, err
	}
	return resp.LastInsertId()
}

func GetProblemJudgeByProblemId(problem_id string) (pj ProblemJudgeType, err error) {
	row := db.QueryRow("SELECT `id`,`judge_case_mode` FROM problem WHERE `problem_id`=?", problem_id)
	var judgeCaseMod string
	err = row.Scan(&pj.Pid, &judgeCaseMod)
	if err != nil {
		return
	}
	// 如果是遇错即止，那么占用资源为1
	if judgeCaseMod == judge.ERGODIC_WITHOUT_ERROR {
		pj.Parallelism = 1
	} else {
		// 否则为样例数量
		pj.Parallelism, err = GetCountProblemCaseByPid(pj.Pid)
	}
	return
}

func GetCountProblemCaseByPid(pid int) (countCase int, err error) {
	row := db.QueryRow("SELECT COUNT(*) FROM problem_case WHERE `pid`=?", pid)
	err = row.Scan(&countCase)
	return
}
