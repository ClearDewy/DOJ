/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package redis

import (
	"doj-go/DataBackup/config"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init(conf *config.Config) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPassword,
	})

}
