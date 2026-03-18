package service

import (
	"fmt"
	"net/http"

	"github.com/Jia-1-svg/crawler/practice/api/basic/config"
	"github.com/Jia-1-svg/crawler/practice/api/handler/request"
	__ "github.com/Jia-1-svg/crawler/practice/api/proto"
	"github.com/Jia-1-svg/crawler/practice/rpc/mq/mainpublish"

	"github.com/gin-gonic/gin"
)

func GoodsAdd(c *gin.Context) {
	var form request.GoodsAdd
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"code": 500,
		})
		return
	}
	add, err := config.UserClient.GoodsAdd(c, &__.GoodsAddReq{
		Title: form.Title,
		Price: float32(form.Price),
		Stock: int64(form.Stock),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "商品发布失败" + err.Error(),
			"code": 500,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "商品发布成功",
		"code": 200,
		"data": add,
	})
	return
}
func GoodsDel(c *gin.Context) {
	var form request.GoodsDel
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"code": 500,
		})
		return
	}
	add, err := config.UserClient.GoodsDel(c, &__.GoodsDelReq{
		GoodsId: int64(form.GoodsId),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "商品删除失败" + err.Error(),
			"code": 500,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "商品删除成功",
		"code": 200,
		"data": add,
	})
	return
}

func GoodsUpdate(c *gin.Context) {
	var form request.GoodsUpdate
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"code": 500,
		})
		return
	}
	add, err := config.UserClient.GoodsUpdate(c, &__.GoodsUpdateReq{
		Title:   form.Title,
		Price:   float32(form.Price),
		Stock:   int64(form.Stock),
		GoodsId: int64(form.GoodsId),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "商品修改失败" + err.Error(),
			"code": 500,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "商品修改成功",
		"code": 200,
		"data": add,
	})
	return
}

func GoodsDetail(c *gin.Context) {
	var form request.GoodsDetail
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"code": 500,
		})
		return
	}
	add, err := config.UserClient.GoodsDetail(c, &__.GoodsDetailReq{
		GoodsId: int64(form.GoodsId),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "商品详情查询失败" + err.Error(),
			"code": 500,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "商品详情查询成功",
		"code": 200,
		"data": add,
	})
	return
}

func OrderAdd(c *gin.Context) {
	var form request.OrderAdd
	// 根据 Content-Type Header 推断使用哪个绑定器。
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"code": 500,
		})
		return
	}
	var list []*__.OrderItem
	for _, adds := range form.OrderItem {
		list = append(list, &__.OrderItem{
			GoodsId:  int64(adds.GoodsId),
			Quantity: int64(adds.Quantity),
		})
	}
	add, err := config.UserClient.OrderAdd(c, &__.OrderAddReq{
		UserId:  int64(form.UserId),
		PayType: int64(form.PayType),
		Item:    list,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "订单添加失败" + err.Error(),
			"code": 500,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "订单添加成功",
		"code": 200,
		"data": add,
	})
	return
}
func NotifyPay(c *gin.Context) {
	c.Request.ParseForm()
	fmt.Println(c.Request.PostForm)
	trade_status := c.Request.PostForm.Get("trade_status")
	out_trade_no := c.Request.PostForm.Get("out_trade_no")
	if trade_status != "TRADE_SUCCESS" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "支付失败",
			"code": 500,
		})
		return
	}
	if out_trade_no == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "订单号不存在",
			"code": 500,
		})
		return
	}
	go func() {
		mainpublish.MainpublishSendMsg("topic", out_trade_no)
	}()
	//pay, err := config.UserClient.OrderNotifyPay(c, &__.OrderNotifyPayReq{
	//	OrderSn: out_trade_no,
	//})
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"msg":  "异步失败",
	//		"code": 500,
	//	})
	//	return
	//}
	c.JSON(200, gin.H{
		"msg":  "异步成功",
		"code": 200,
		//"data": pay,
	})
	return
}
