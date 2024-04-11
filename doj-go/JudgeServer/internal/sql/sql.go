/**
 * @ Author: ClearDewy
 * @ Desc: 初始化
 **/
package sql

import (
	"context"
	"database/sql"
	"doj-go/JudgeServer/config"
	"fmt"
	"github.com/ClearDewy/go-pkg/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func Init() (func() error, func(ctx context.Context) error) {
	logrus.Info(fmt.Sprintf("%s:%s@tcp(%s)/doj", config.Conf.MysqlUsername, config.Conf.MysqlPassword, config.Conf.MysqlAddr))
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/doj", config.Conf.MysqlUsername, config.Conf.MysqlPassword, config.Conf.MysqlAddr))
	if err != nil {
		logrus.ErrorM(err, "数据库连接失败")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		logrus.ErrorM(err, "数据库Ping失败")
		panic(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return nil, stop
}

func stop(ctx context.Context) error {
	return db.Close()
}
