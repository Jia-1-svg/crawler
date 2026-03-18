package mainpublish

import (
	"zhongyao/aa/crawler/practice/rpc/mq/RabbitMQ"
)

func MainpublishSendMsg(topic string, orderSn string) {
	kutengOne := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")

	kutengOne.SendMsg(topic, orderSn)

}
