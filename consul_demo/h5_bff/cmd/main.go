package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"consul_demo/h5_bff/client"
	"consul_demo/h5_bff/config"
	"consul_demo/h5_bff/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.DefaultConfig()

	// 设置 gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化 Order 客户端（通过 Consul 发现服务）
	orderClient, err := client.NewOrderClient(cfg.Consul.Address, cfg.Order.ServiceName)
	if err != nil {
		log.Fatalf("Failed to create order client: %v", err)
	}
	log.Printf("Order client initialized, service discovered: %s", cfg.Order.ServiceName)

	// 创建订单处理器
	orderHandler := handler.NewOrderHandler(orderClient)

	// 初始化 Gin 路由
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册 HTTP 路由
	api := r.Group("/")
	{
		api.POST("/orders", orderHandler.CreateOrderHandler)
		api.GET("/orders/:id", orderHandler.GetOrderHandler)
		api.PUT("/orders/:id", orderHandler.UpdateOrderHandler)
		api.DELETE("/orders/:id", orderHandler.DeleteOrderHandler)
		api.GET("/orders", orderHandler.ListOrdersHandler)
	}

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("H5 BFF server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
