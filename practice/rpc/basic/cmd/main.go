package main

import (
	"flag"
	"log"
	"net"
	"zhongyao/aa/crawler/practice/rpc/RabbitMQ"
	"zhongyao/aa/crawler/practice/rpc/handler"
	__ "zhongyao/aa/crawler/practice/rpc/proto"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//
	//services, err := init1.GetServiceWithLoadBalancer(config.Config.Consul.ServiceName)
	//if err != nil {
	//	log.Printf("获取用户服务失败: %v", err)
	//} else {
	//	log.Printf("获取到用户服务: %s, 地址: %s:%d", services.Service, services.Address, services.Port)
	//}
	//
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})
	__.RegisterUserServiceServer(s, &handler.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	//
	//go func() {
	//	if err := s.Serve(lis); err != nil {
	//		log.Fatalf("failed to serve: %v", err)
	//	}
	//}()
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//log.Println("正在关闭服务...")
	//if err := init1.ConsulShutdown(); err != nil {
	//	log.Printf("Consul注销失败: %v", err)
	//}
	//_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//s.GracefulStop()
	//log.Println("服务已关闭")
	RabbitMQ.ConsumeStockDeduct()
}
