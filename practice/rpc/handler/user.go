package handler

import (
	"context"
	"errors"
	"github.com/Jia-1-svg/crawler/practice/rpc/RabbitMQ"
	"github.com/Jia-1-svg/crawler/practice/rpc/basic/config"
	"github.com/Jia-1-svg/crawler/practice/rpc/model"
	"github.com/Jia-1-svg/crawler/practice/rpc/pkg"
	__ "github.com/Jia-1-svg/crawler/practice/rpc/proto"
	"strconv"
)

type Server struct {
	__.UnimplementedUserServiceServer
	//pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) UserRegister(_ context.Context, in *__.UserRegisterReq) (*__.UserRegisterResp, error) {
	user := model.User{
		Nickname: in.Nickname,
		Password: in.Password,
	}
	err := user.UserAdd(config.DB)
	if err != nil {
		return nil, errors.New("用户注册失败")
	}
	return &__.UserRegisterResp{
		UserId: int64(user.ID),
	}, nil
}
func (s *Server) OrderAdd(_ context.Context, in *__.OrderAddReq) (*__.OrderAddResp, error) {
	orderSn := pkg.CreateOrderSn()
	total := 0.0
	var orderItem []*model.OrderItem
	for _, item := range in.Item {
		var goods model.Goods
		err := goods.FindGoodsById(config.DB, item.GoodsId)
		if err != nil {
			return nil, errors.New("商品不存在")
		}
		total += goods.Price * float64(item.Quantity)
		orderItem = append(orderItem, &model.OrderItem{
			GoodsId:  int(item.GoodsId),
			Quantity: int(item.Quantity),
		})
	}
	order := model.Order{
		UserId:  int(in.UserId),
		PayType: int(in.PayType),
		OrderSn: orderSn,
		Status:  0,
		Total:   total,
	}
	err := order.OrderAdd(config.DB)
	if err != nil {
		return nil, errors.New("订单添加失败")
	}
	for i, _ := range orderItem {
		orderItem[i].OrderId = int(order.ID)
	}
	err = order.OrderItemAdd(config.DB, orderItem)
	if err != nil {
		return nil, errors.New("订单详情查询失败")
	}
	pay := pkg.AliPay(orderSn, total)

	return &__.OrderAddResp{
		OrderId: int64(order.ID),
		AliPay:  pay,
		Total:   float32(total),
	}, nil
}
func (s *Server) GoodsAdd(_ context.Context, in *__.GoodsAddReq) (*__.GoodsAddResp, error) {
	goods := model.Goods{
		Title: in.Title,
		Price: float64(in.Price),
		Stock: int(in.Stock),
	}
	err := goods.GoodsAdd(config.DB)
	if err != nil {
		return nil, errors.New("商品创建失败")
	}
	return &__.GoodsAddResp{
		GoodsId: int64(goods.ID),
	}, nil
}
func (s *Server) GoodsDel(_ context.Context, in *__.GoodsDelReq) (*__.GoodsDelResp, error) {
	var goods model.Goods
	err := goods.GoodsDel(config.DB, in.GoodsId)
	if err != nil {
		return nil, errors.New("商品删除失败")
	}
	return &__.GoodsDelResp{
		Success: true,
	}, nil
}
func (s *Server) GoodsUpdate(_ context.Context, in *__.GoodsUpdateReq) (*__.GoodsUpdateResp, error) {
	var goods model.Goods
	err := goods.FindGoodsById(config.DB, in.GoodsId)
	if err != nil {
		return nil, errors.New("商品查询失败")
	}
	goods.Price = float64(in.Price)
	goods.Title = in.Title
	goods.Stock = int(in.Stock)
	err = goods.UpdateGoods(config.DB)
	if err != nil {
		return nil, errors.New("商品信息修改失败")
	}
	return &__.GoodsUpdateResp{
		Success: true,
	}, nil
}
func (s *Server) GoodsDetail(_ context.Context, in *__.GoodsDetailReq) (*__.GoodsDetailResp, error) {
	var goods model.Goods
	err := goods.FindGoodsById(config.DB, in.GoodsId)
	if err != nil {
		return nil, errors.New("商品查询失败")
	}
	return &__.GoodsDetailResp{
		Title: goods.Title,
		Price: float32(goods.Price),
		Stock: int64(goods.Stock),
	}, nil
}
func (s *Server) OrderNotifyPay(_ context.Context, in *__.OrderNotifyPayReq) (*__.OrderNotifyPayResp, error) {
	var order model.Order
	//go func() {
	//	sendmsg.SendMsg(in.OrderSn)
	//}()
	err := order.FindOrderByOrderSn(config.DB, in.OrderSn)
	if err != nil {
		return nil, errors.New("订单查询失败")
	}
	order.Status = 1
	err = order.UpdateOrderStatus(config.DB)
	if err != nil {
		return nil, errors.New("订单状态更新失败")
	}
	var orderItem []model.OrderItem
	err = config.DB.Debug().Where("order_id =?", order.ID).Find(&orderItem).Error
	if err != nil {
		return nil, errors.New("订单详情查询失败")
	}
	for _, item := range orderItem {
		var goods model.Goods
		RabbitMQ.SendStockDeductMsg(strconv.FormatInt(int64(item.GoodsId), 10), item.Quantity)
		err = goods.FindGoodsById(config.DB, int64(item.GoodsId))
		if err != nil {
			return nil, errors.New("商品查询失败")
		}
		goods.Stock -= item.Quantity
		err := goods.UpdateGoods(config.DB)
		if err != nil {
			return nil, errors.New("库存修改失败")
		}
	}
	return &__.OrderNotifyPayResp{
		Success: true,
	}, nil
}
