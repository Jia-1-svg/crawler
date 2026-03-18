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
	"zhongyao/aa/crawler/practice/rpc/handler"
	__ "zhongyao/aa/crawler/practice/rpc/proto"

	"google.golang.org/grpc"

	"github.com/gospacex/gospacex/core/storage/registry/consul"
)

const (
	serviceName    = "user-service"
	serviceAddress = "127.0.0.1"
	servicePort    = 9090
	consulAddr     = "115.190.88.2:8500"
)

func main() {
	//建立连接返回操作consul的句柄
	client, err := consul.NewClient(consulAddr)
	if err != nil {
		log.Fatalf("Failed to create consul client: %v", err)
	}
	//可通过两种方式写活入参:
	//1,flag:https://github.com/spf13/pflag
	//2.环境变量: os.Getenv("SERVICE_ID")
	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, serviceAddress, servicePort)
	//标签
	tags := []string{"grpc", "v1", "user"}

	//服务注册
	err = client.RegisterServiceWithTTL(serviceID, serviceName, serviceAddress, servicePort, tags, "10s")
	if err != nil {
		log.Fatalf("Failed to register service to consul: %v", err)
	}
	log.Printf("Service registered to consul: %s", serviceID)

	//监听tcp端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serviceAddress, servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	//初始化grpc服务
	grpcServer := grpc.NewServer()
	//注册grpc服务
	__.RegisterUserServiceServer(grpcServer, &handler.Server{})

	//定义上下文清除
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//通过定时器进行健康检查
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := client.PassTTL(serviceID, "service is healthy"); err != nil {
					log.Printf("Failed to update TTL: %v", err)
				} else {
					log.Println("TTL heartbeat sent")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	//启动grpc
	go func() {
		log.Printf("gRPC server starting on %s:%d", serviceAddress, servicePort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Failed to serve gRPC: %v", err)
		}
	}()

	//接收信号进行优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	cancel()
	grpcServer.GracefulStop()
	//注销微服
	if err := client.DeregisterService(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}
	log.Println("Service deregistered from consul")
}
