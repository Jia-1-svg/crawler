package main

import (
	"context"
	"log"
	"time"
	config1 "zhongyao/aa/crawler/practice/api/basic/config"
	__ "zhongyao/aa/crawler/practice/api/proto"
	"zhongyao/aa/crawler/practice/rpc/basic/config"
	"zhongyao/aa/crawler/practice/rpc/mq/RabbitMQ"
)

func main() {
	kutengOne := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "#")
	kutengOne.SubsribeMsg("topic", func(msg string) {
		val := config.Rdb.SetNX(config.Ctx, "orderSn", 1, time.Second*10).Val()
		if !val {
			log.Println("11")
		}

		config1.UserClient.OrderNotifyPay(context.Background(), &__.OrderNotifyPayReq{
			OrderSn: msg,
		})

		return
	})
}
