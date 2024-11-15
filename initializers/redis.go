package initializers

import (
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "litstore_redis:6379",
	})
}
