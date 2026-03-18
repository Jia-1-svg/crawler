package sendmsg

import (
	RabbitMQ "zhongyao/aa/crawler/practice/rpc/rabbitmqs/rabbitmq"
)

func SendMsg(orderSn string) {
	kutengOne := RabbitMQ.NewRabbitMQTopic("topic", "kuteng.topic.one")
	kutengOne.PublishTopic("topic", orderSn)
}
