/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package config

import (
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	ROOT_PATH = "/doj/judge"
)

var Conf = &Config{}

type Config struct {
	Parallelism      int    `flagUsage:"control the # of concurrency execution (default equal to number of cpu)"`
	MysqlAddr        string `default:"127.0.0.1:3306"`
	MysqlUsername    string `default:"root"`
	MysqlPassword    string `default:"doj123456"`
	EtcdAddr         string `flagUsage:"etcd主机" default:"127.0.0.1:2379"`
	EtcdRootPassword string `flagUsage:"etcd root 的密码" default:"doj123456"`
	JudgeServerName  string `flagUsage:"显示在后端的判题机的名称" default:"judge-server-alone"`
	JudgeServerAddr  string `flagUsage:"判题服务器的公网" default:"127.0.0.1:8888"`
	OutputLimit      int64  `flagUsage:"最大输出限制" default:"256<<20"`
}

func (c *Config) LoadEnvDefault() (err error) {
	err = loadEnvDefault(c)
	if c.Parallelism <= 0 {
		c.Parallelism = runtime.NumCPU()
	}
	return
}

func loadEnvDefault(v interface{}) error {
	elem := reflect.ValueOf(v).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		typeField := elem.Type().Field(i)
		var err error
		envName := camelCaseToEnvVar(typeField.Name)
		if envVal, ok := os.LookupEnv(envName); field.CanSet() {
			if ok {
				err = setFieldWithValue(field, envVal)
			} else {
				if tagValue := typeField.Tag.Get("default"); tagValue != "" {
					err = setFieldWithValue(field, tagValue)
				}
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}
func camelCaseToEnvVar(name string) string {
	// 正则表达式匹配大写字母
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	words := re.FindAllString(name, -1)
	for i := range words {
		words[i] = strings.ToUpper(words[i])
	}
	return strings.Join(words, "_")
}
func setFieldWithValue(field reflect.Value, value string) error {
	if field.CanSet() {
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int64:
			// 首先检查字段类型是否为 time.Duration
			if field.Type() == reflect.TypeOf(time.Duration(0)) {
				if duration, err := time.ParseDuration(value); err == nil {
					field.Set(reflect.ValueOf(duration))
				} else {
					return err
				}
			} else {
				// 处理普通的整数
				if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
					field.SetInt(intValue)
				} else {
					return err
				}
			}
		case reflect.Bool:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				field.SetBool(boolValue)
			} else {
				return err
			}

		default:
			field.Set(reflect.ValueOf(value))
		}
	}
	return nil
}
