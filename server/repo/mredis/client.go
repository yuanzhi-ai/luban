package mredis

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

const (
	EmptyKeyErr = redis.Nil
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("WANXIANG_REDIS_ADDR"),
		Password: os.Getenv("WANXIANG_REDIS_PSWD"),
		DB:       0, // use default DB
	})
}

// 获取一个rdb的客户端
func GetRdbClient() *redis.Client {
	return rdb
}
