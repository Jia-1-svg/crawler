package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"consul_demo/order_srv/config"
	"consul_demo/order_srv/model"
	"consul_demo/order_srv/service"
	"consul_demo/pkg/consul"
	pb "consul_demo/proto/order"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg := config.DefaultConfig()

	// 初始化数据库连接
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// 初始化 Consul 客户端
	consulClient, err := consul.NewClient(cfg.Consul.Address)
	if err != nil {
		log.Fatalf("Failed to create consul client: %v", err)
	}
	log.Printf("Consul client connected to %s", cfg.Consul.Address)

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册订单服务
	orderService := service.NewOrderService(db)
	pb.RegisterOrderServiceServer(grpcServer, orderService)

	// 启动 gRPC 服务器
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 注册服务到 Consul（使用 TTL 健康检查）
	err = consulClient.RegisterService(
		cfg.Service.ServiceID,
		cfg.Service.ServiceName,
		cfg.Service.Host,
		cfg.Service.Port,
		[]string{"order", "grpc"},
	)
	if err != nil {
		log.Fatalf("Failed to register service to consul: %v", err)
	}
	log.Printf("Service registered to Consul: %s (ID: %s)", cfg.Service.ServiceName, cfg.Service.ServiceID)

	// 启动 TTL 健康检查更新协程
	go func() {
		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := consulClient.PassTTL(cfg.Service.ServiceID, "service healthy"); err != nil {
					log.Printf("Failed to update TTL: %v", err)
				}
			}
		}
	}()

	// 优雅关闭处理
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down server...")

		// 从 Consul 注销服务
		if err := consulClient.DeregisterService(cfg.Service.ServiceID); err != nil {
			log.Printf("Failed to deregister service: %v", err)
		} else {
			log.Println("Service deregistered from Consul")
		}

		// 优雅关闭 gRPC 服务器
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			log.Println("Server stopped gracefully")
		case <-ctx.Done():
			log.Println("Force stopping server")
			grpcServer.Stop()
		}

		os.Exit(0)
	}()

	log.Printf("Order service started on port %d", cfg.Service.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// initDB 初始化数据库连接
func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4\u0026parseTime=True\u0026loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&model.Order{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
