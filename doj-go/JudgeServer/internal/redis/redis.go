/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package redis

import (
	"context"
	"doj-go/JudgeServer/config"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init() (func() error, func(ctx context.Context) error) {

	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisAddr,
		Password: config.Conf.RedisPassword,
	})
	return nil, nil
}
