package main

import (
	"zhongyao/aa/crawler/practice/rpc/RabbitMQ"
	_ "zhongyao/aa/crawler/practice/rpc/basic/init1"
)

func main() {
	RabbitMQ.ConsumeStockDeduct()
}
