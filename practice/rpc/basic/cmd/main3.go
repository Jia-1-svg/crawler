package main

import (
	"github.com/Jia-1-svg/crawler/practice/rpc/RabbitMQ"
	_ "github.com/Jia-1-svg/crawler/practice/rpc/basic/init1"
)

func main() {
	RabbitMQ.SendStockDeductMsg("苹果", 1)
}
