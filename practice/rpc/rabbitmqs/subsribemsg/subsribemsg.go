package subsribemsg

import (
	"log"
	"time"
	"zhongyao/aa/crawler/practice/rpc/basic/config"
	RabbitMQ "zhongyao/aa/crawler/practice/rpc/rabbitmqs/rabbitmq"
)

func main() {
	kutengOne := RabbitMQ.NewRabbitMQTopic("topic", "#")
	kutengOne.RecieveTopic("topic", func(msg string) {
		key := "order_sn:" + msg
		val := config.Rdb.SetNX(config.Ctx, key, 1, 10*time.Second).Val()
		if !val {
			log.Println()
			return
		}

	})
}
