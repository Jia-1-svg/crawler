package init1

import (
	"fmt"
	"github.com/Jia-1-svg/crawler/practice/rpc/basic/config"

	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	redisConfig := config.Config.Redis
	config.Rdb = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.Database, // use default DB
	})

	err := config.Rdb.Set(config.Ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis连接成功")
}
