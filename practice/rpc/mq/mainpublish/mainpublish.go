package mainpublish

import rabbitMQ "github.com/Jia-1-svg/crawler/practice/rpc/mq/RabbitMQ"

func MainpublishSendMsg(topic string, orderSn string) {
	kutengOne := rabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	kutengOne.SendMsg(topic, orderSn)

}
