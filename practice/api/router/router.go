package router

import (
	"zhongyao/aa/crawler/practice/api/handler/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.POST("/goods/add", service.GoodsAdd)
	r.POST("/goods/del", service.GoodsDel)
	r.POST("/goods/update", service.GoodsUpdate)
	r.GET("/goods/detail", service.GoodsDetail)
	r.POST("/order/add", service.OrderAdd)
	r.POST("/notify/pay", service.NotifyPay)
	return r
}
