/**
 * @ Author: ClearDewy
 * @ Desc: 用于加载后端启动时的配置
 **/
package config

import (
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BackendServerPort string `default:"8080"`
	MysqlAddr         string `default:"127.0.0.1:3306"`
	MysqlUsername     string `default:"root"`
	MysqlPassword     string `default:"doj123456"`
	RedisAddr         string `default:"127.0.0.1:6379"`
	RedisPassword     string `default:"doj123456"`
	EtcdAddr          string `default:"127.0.0.1:2379"`
	EtcdRootPassword  string `default:"doj123456"`
}

func (c *Config) LoadEnvDefault() error {
	return loadEnvDefault(c)
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
