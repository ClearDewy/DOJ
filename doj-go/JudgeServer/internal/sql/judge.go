/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sql

import (
	"time"
)

type JudgeInfoType struct {
	Lid  int
	Code string
}

type JudgeCaseStatusType struct {
	Jid           int
	ProblemCaseId int
	Status        int
	Message       string
	Time          uint64
	Memory        uint64
}

type JudgeStatusType struct {
	Jid     int
	Status  int
	Time    uint64
	Memory  uint64
	Message string
}

func UpdateJudgeStatus(judgeStatus *JudgeStatusType) (err error) {
	_, err = db.Exec("UPDATE judge SET `status`= ?,`time`=?,`memory`=?,`message`=? WHERE id=?",
		judgeStatus.Status, judgeStatus.Time/uint64(time.Millisecond), judgeStatus.Memory>>20, judgeStatus.Message, judgeStatus.Jid)
	return
}

func GetJudgeCode(jid int) (judgeInfo JudgeInfoType, err error) {
	row := db.QueryRow("SELECT `lid`,`code` FROM judge WHERE id=?", jid)
	err = row.Scan(&judgeInfo.Lid, &judgeInfo.Code)
	return
}

func AddOrUpdateJudgeCaseStatus(jcs *JudgeCaseStatusType) (err error) {
	_, err = db.Exec(`
    INSERT INTO judge_case(jid, problem_case_id, status, message, time, memory)
    VALUES(?,?,?,?,?,?)
    ON DUPLICATE KEY UPDATE status=?, message=?, time=?, memory=?`,
		jcs.Jid, jcs.ProblemCaseId, jcs.Status, jcs.Memory, jcs.Time/uint64(time.Millisecond), jcs.Memory>>20,
		jcs.Status, jcs.Memory, jcs.Time/uint64(time.Millisecond), jcs.Memory>>20)
	return
}
