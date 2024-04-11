/**
 * @ Author: ClearDewy
 * @ Desc: 初始化
 **/
package sql

import (
	"database/sql"
	"doj-go/DataBackup/config"
	"fmt"
	"github.com/ClearDewy/go-pkg/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func Init(conf *config.Config) {
	logrus.Info(fmt.Sprintf("%s:%s@tcp(%s)/doj", conf.MysqlUsername, conf.MysqlPassword, conf.MysqlAddr))
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/doj", conf.MysqlUsername, conf.MysqlPassword, conf.MysqlAddr))
	logrus.ErrorM(err, "数据库连接失败")
	err = db.Ping()
	logrus.ErrorM(err, "数据库Ping失败")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

type String string
type Bool bool
type Int int

func (s *String) Scan(value interface{}) error {
	if value == nil {
		*s = ""
	} else if bv, ok := value.([]byte); ok {
		// 如果value是[]byte类型，则将其转换为string
		*s = String(bv)
	} else {
		// 您可以根据需要在这里添加更多的类型检查
		return fmt.Errorf("cannot scan type %T into String", value)
	}
	return nil
}
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		*b = false
	} else if iv, ok := value.(int64); ok {
		// 如果value是int64类型，则将其转换为Bool
		*b = iv != 0
	} else {
		return fmt.Errorf("cannot scan type %T into Bool", value)
	}
	return nil
}
func (i *Int) Scan(value interface{}) error {
	if value == nil {
		*i = 0
	} else if iv, ok := value.(int64); ok {
		// 如果value是int64类型，则将其转换为Int
		*i = Int(iv)
	} else {
		return fmt.Errorf("cannot scan type %T into Int", value)
	}
	return nil
}
