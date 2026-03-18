package main

import (
	"context"
	config1 "github.com/Jia-1-svg/crawler/practice/api/basic/config"
	_ "github.com/Jia-1-svg/crawler/practice/api/basic/init"
	__ "github.com/Jia-1-svg/crawler/practice/api/proto"
	"github.com/Jia-1-svg/crawler/practice/rpc/basic/config"
	_ "github.com/Jia-1-svg/crawler/practice/rpc/basic/init1"
	rabbitMQ "github.com/Jia-1-svg/crawler/practice/rpc/mq/RabbitMQ"
	"log"
	"time"
)

func main() {
	kutengOne := rabbitMQ.NewRabbitMQTopic("exKutengTopic", "#")
	kutengOne.SubsribeMsg("topic", func(msg string) {
		val := config.Rdb.SetNX(context.Background(), msg, 1, time.Minute*10).Val()
		if !val {
			log.Println("错误")
			return
		}
		_, err := config1.UserClient.OrderNotifyPay(context.Background(), &__.OrderNotifyPayReq{
			OrderSn: msg,
		})
		if err != nil {
			log.Println("错误22")
		}
	})
}
