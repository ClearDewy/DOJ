/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sql

import (
	"fmt"
	"github.com/ClearDewy/go-pkg/logrus"
	dSql "github.com/ClearDewy/go-pkg/sql"
	"strings"
)

type TagType struct {
	Id    int         `json:"id"`
	Name  dSql.String `json:"name"`
	Color dSql.String `json:"color"`
}

type ProblemListType struct {
	Problem_id dSql.String `json:"problem_id"`
	Title      dSql.String `json:"title"`
	Difficulty int         `json:"difficulty"`
	Tags       []*TagType  `json:"tags"`
	MyStatus   int         `json:"myStatus"`
	Total      int         `json:"total"`
	Ac         int         `json:"ac"`
	Wa         int         `json:"wa"`
	Tle        int         `json:"tle"`
	Mle        int         `json:"mle"`
	Re         int         `json:"re"`
	Pe         int         `json:"pe"`
	Ce         int         `json:"ce"`
	Se         int         `json:"se"`
}

type ProblemLanguageType struct {
	Lid          int         `json:"lid"`
	Content_type dSql.String `json:"content_type"`
	Name         dSql.String `json:"name"`
	Time_limit   int         `json:"time_limit"`
	Memory_limit int         `json:"memory_limit"`
}

type ProblemDetailType struct {
	Title  dSql.String `json:"title"`
	Author dSql.String `json:"author"`

	Description dSql.String            `json:"description"`
	Input       dSql.String            `json:"input"`
	Output      dSql.String            `json:"output"`
	Examples    dSql.String            `json:"examples"`
	Oj          dSql.String            `json:"oj"`
	Hint        dSql.String            `json:"hint"`
	Source      dSql.String            `json:"source"`
	Tags        []*TagType             `json:"tags"`
	Languages   []*ProblemLanguageType `json:"languages"`
}

func GetProblemList(oj, difficulty, keyword string, tags []string, limit, currentPage int) (problemList []*ProblemListType, err error) {
	wheres := make([]string, 0, 3)
	if difficulty != "" && !strings.EqualFold(difficulty, "all") {
		wheres = append(wheres, fmt.Sprintf("`difficulty`='%s'", difficulty))
	}
	if oj != "" && !strings.EqualFold(oj, "all") {
		wheres = append(wheres, fmt.Sprintf("`oj`='%s'", oj))
	}
	if keyword != "" {
		wheres = append(wheres, fmt.Sprintf("`title`LIKE '%%%s%%'", keyword))
	}

	where := ""
	if len(wheres) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(wheres, "&&"))
	}

	sql := "SELECT `problem`.id,`problem_id`,`title`,`difficulty` FROM problem"
	if tags != nil && len(tags) > 0 {
		sql = fmt.Sprintf("%s JOIN (SELECT `pid` FROM problem_tag WHERE `tid` IN (%s) GROUP BY `pid` HAVING COUNT(*)=%d) ON `pid`=problem.id",
			sql, strings.Join(tags, ","), len(tags))
	}
	sql = fmt.Sprintf("%s %s LIMIT %d OFFSET %d", sql, where, limit, (currentPage-1)*limit)
	logrus.Info(sql)
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		pli := ProblemListType{}
		pid := 0
		err = rows.Scan(&pid, &pli.Problem_id, &pli.Title, &pli.Difficulty)
		if err != nil {
			return
		}
		pli.Tags, err = GetProblemTags(pid)
		problemList = append(problemList, &pli)
	}
	return
}

func GetProblemTags(pid int) (tags []*TagType, err error) {
	rows, err := db.Query("SELECT tag.id,`name`,`color` FROM tag JOIN problem_tag on tag.id = problem_tag.tid WHERE problem_tag.pid=?", pid)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		tag := TagType{}
		err = rows.Scan(&tag.Id, &tag.Name, &tag.Color)
		if err != nil {
			return
		}
		tags = append(tags, &tag)
	}
	return
}

func GetProblemLanguage(pid int) (problem_languages []*ProblemLanguageType, err error) {
	rows, err := db.Query("SELECT `lid`,`name`,`content_type`,`time_limit`,`memory_limit` FROM problem_language JOIN language ON problem_language.lid=language.id WHERE pid=?", pid)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		lang := &ProblemLanguageType{}
		err = rows.Scan(&lang.Lid, &lang.Name, &lang.Content_type, &lang.Time_limit, &lang.Memory_limit)
		if err != nil {
			return
		}
		problem_languages = append(problem_languages, lang)
	}
	return
}

func GetProblemDetail(problem_id string) (pstmt *ProblemDetailType, err error) {
	pstmt = &ProblemDetailType{}
	pid := 0
	row := db.QueryRow("SELECT `id`,`title`,IFNULL(`author`,''),IFNULL(`description`,''),IFNULL(`input`,''),IFNULL(`output`,''),IFNULL(`examples`,''),`oj`,IFNULL(`hint`,''),IFNULL(`source`,'') FROM problem WHERE `problem_id`= ? ", problem_id)
	err = row.Scan(&pid, &pstmt.Title, &pstmt.Author,
		&pstmt.Description, &pstmt.Input, &pstmt.Output, &pstmt.Examples, &pstmt.Oj, &pstmt.Hint, &pstmt.Source)
	if err != nil {
		return
	}
	pstmt.Languages, err = GetProblemLanguage(pid)
	if err != nil {
		return
	}
	pstmt.Tags, err = GetProblemTags(pid)
	if err != nil {
		return
	}
	return
}

func GetPidByProblemId(problem_id string) (pid int, err error) {
	row := db.QueryRow("SELECT `id` from problem WHERE problem_id=?", problem_id)
	err = row.Scan(&pid)
	return
}
