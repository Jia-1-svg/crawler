package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"consul_demo/order_srv/model"
	pb "consul_demo/proto/order"

	"gorm.io/gorm"
)

// OrderService 订单服务实现
type OrderService struct {
	pb.UnimplementedOrderServiceServer
	db *gorm.DB
}

// NewOrderService 创建订单服务
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// generateOrderNo 生成订单号: yyyyMMdd + 6位随机数
func generateOrderNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	randomNum := rand.Intn(900000) + 100000
	return fmt.Sprintf("%s%d", dateStr, randomNum)
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := &model.Order{
		OrderNo:    generateOrderNo(),
		UserID:     req.UserId,
		TotalPrice: req.TotalPrice,
		Status:     model.STATUS_CREATE,
		CreateTime: model.GetCurrentTime(),
		UpdateTime: model.GetCurrentTime(),
	}

	if err := s.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &pb.CreateOrderResponse{
		OrderId: order.ID,
		OrderNo: order.OrderNo,
	}, nil
}

// GetOrder 根据ID查询订单
func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	var order model.Order
	if err := s.db.WithContext(ctx).First(&order, req.OrderId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found: %d", req.OrderId)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return convertToProtoOrder(&order), nil
}

// UpdateOrder 更新订单
func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	updates := map[string]interface{}{
		"update_time": model.GetCurrentTime(),
	}

	if req.TotalPrice > 0 {
		updates["total_price"] = req.TotalPrice
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	result := s.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", req.OrderId).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return &pb.UpdateOrderResponse{Success: false}, fmt.Errorf("order not found: %d", req.OrderId)
	}

	return &pb.UpdateOrderResponse{Success: true}, nil
}

// DeleteOrder 删除订单
func (s *OrderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	result := s.db.WithContext(ctx).Delete(&model.Order{}, req.OrderId)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to delete order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return &pb.DeleteOrderResponse{Success: false}, fmt.Errorf("order not found: %d", req.OrderId)
	}

	return &pb.DeleteOrderResponse{Success: true}, nil
}

// ListOrders 分页查询订单列表
func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	var orders []model.Order
	var total int64

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	query := s.db.WithContext(ctx).Model(&model.Order{})
	if req.UserId > 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Limit(int(pageSize)).Offset(int(offset)).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = convertToProtoOrder(&order)
	}

	return &pb.ListOrdersResponse{
		Orders: pbOrders,
		Total:  int32(total),
	}, nil
}

// convertToProtoOrder 将 model.Order 转换为 pb.Order
func convertToProtoOrder(order *model.Order) *pb.Order {
	return &pb.Order{
		Id:         order.ID,
		OrderNo:    order.OrderNo,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreateTime: order.CreateTime,
		UpdateTime: order.UpdateTime,
	}
}
