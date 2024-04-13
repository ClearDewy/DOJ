/**
 * @ Author: ClearDewy
 * @ Desc: 用于加载后端启动时的配置
 **/
package config

import (
	"github.com/koding/multiconfig"
	"time"
)

var Conf = &Config{}

type Config struct {
	BackendServerPort string `default:"8080"`
	MysqlAddr         string `default:"127.0.0.1:3306"`
	MysqlUsername     string `default:"root"`
	MysqlPassword     string `default:"doj123456"`
	RedisAddr         string `default:"127.0.0.1:6379"`
	RedisPassword     string `default:"doj123456"`
	EtcdAddr          string `default:"127.0.0.1:2379"`
	EtcdRootPassword  string `default:"doj123456"`

	JudgeServerExpireTime time.Duration `default:"1h"`
}

func (c *Config) LoadEnv() error {
	cl := multiconfig.MultiLoader(
		&multiconfig.EnvironmentLoader{
			Prefix:    "doj",
			CamelCase: true,
		},
		&multiconfig.TagLoader{},
		&multiconfig.FlagLoader{
			Prefix:    "doj",
			CamelCase: true,
		},
	)
	return cl.Load(c)
}
