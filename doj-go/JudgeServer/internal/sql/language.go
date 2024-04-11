/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sql

import "time"

type LanguageLimitType struct {
	TimeLimit   uint64
	MemoryLimit uint64
	StackLimit  uint64
}

func GetProblemLanguageLimit(pid, lid int) (limit LanguageLimitType, err error) {
	row := db.QueryRow("SELECT `time_limit`,`memory_limit`,`stack_limit` FROM problem_language WHERE `pid`= ? AND `lid`= ? ", pid, lid)
	err = row.Scan(&limit.TimeLimit, &limit.MemoryLimit, &limit.StackLimit)
	// 转换为ms
	limit.TimeLimit *= uint64(time.Millisecond)
	// 转换为mb
	limit.MemoryLimit <<= 20
	// 转换为mb
	limit.StackLimit <<= 20
	return
}
