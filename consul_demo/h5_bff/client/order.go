package client

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"consul_demo/pkg/consul"
	pb "consul_demo/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrderClient Order 服务 gRPC 客户端
type OrderClient struct {
	clients []pb.OrderServiceClient
	consul  *consul.Client
	mu      sync.RWMutex
	index   int
}

// NewOrderClient 通过 Consul 发现服务并创建 gRPC 客户端（支持负载均衡）
func NewOrderClient(consulAddr, serviceName string) (*OrderClient, error) {
	// 创建 Consul 客户端
	consulClient, err := consul.NewClient(consulAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	// 发现服务
	services, err := consulClient.DiscoverService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	// 创建多个 gRPC 连接（负载均衡）
	var clients []pb.OrderServiceClient
	for _, service := range services {
		target := fmt.Sprintf("%s:%d", service.Address, service.Port)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		conn, err := grpc.DialContext(ctx, target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()

		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", target, err)
			continue
		}

		clients = append(clients, pb.NewOrderServiceClient(conn))
		fmt.Printf("Connected to order service: %s\n", target)
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("failed to connect to any order service instance")
	}

	fmt.Printf("Successfully connected to %d order service instances\n", len(clients))
	return &OrderClient{
		clients: clients,
		consul:  consulClient,
		index:   0,
	}, nil
}

// getClient 轮询获取客户端（简单负载均衡）
func (c *OrderClient) getClient() pb.OrderServiceClient {
	c.mu.Lock()
	defer c.mu.Unlock()

	client := c.clients[c.index%len(c.clients)]
	c.index++
	return client
}

// getClientRandom 随机获取客户端（随机负载均衡）
func (c *OrderClient) getClientRandom() pb.OrderServiceClient {
	c.mu.RLock()
	defer c.mu.RUnlock()

	index := rand.Intn(len(c.clients))
	return c.clients[index]
}

// CreateOrder 创建订单
func (c *OrderClient) CreateOrder(ctx context.Context, userID int64, totalPrice float64) (orderID int64, orderNo string, err error) {
	req := &pb.CreateOrderRequest{
		UserId:     userID,
		TotalPrice: totalPrice,
	}

	client := c.getClient()
	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create order: %w", err)
	}

	return resp.OrderId, resp.OrderNo, nil
}

// GetOrder 获取订单详情
func (c *OrderClient) GetOrder(ctx context.Context, orderID int64) (*pb.Order, error) {
	req := &pb.GetOrderRequest{
		OrderId: orderID,
	}

	client := c.getClient()
	resp, err := client.GetOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", orderID, err)
	}

	return resp, nil
}

// UpdateOrder 更新订单
func (c *OrderClient) UpdateOrder(ctx context.Context, orderID int64, totalPrice float64, status int32) error {
	req := &pb.UpdateOrderRequest{
		OrderId:    orderID,
		TotalPrice: totalPrice,
		Status:     status,
	}

	client := c.getClient()
	resp, err := client.UpdateOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update order %d: %w", orderID, err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to update order %d: service returned failure", orderID)
	}

	return nil
}

// DeleteOrder 删除订单
func (c *OrderClient) DeleteOrder(ctx context.Context, orderID int64) error {
	req := &pb.DeleteOrderRequest{
		OrderId: orderID,
	}

	client := c.getClient()
	resp, err := client.DeleteOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete order %d: %w", orderID, err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to delete order %d: service returned failure", orderID)
	}

	return nil
}

// ListOrders 获取订单列表
func (c *OrderClient) ListOrders(ctx context.Context, userID int64, page, pageSize int32) ([]*pb.Order, int32, error) {
	req := &pb.ListOrdersRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	client := c.getClient()
	resp, err := client.ListOrders(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	return resp.Orders, resp.Total, nil
}
