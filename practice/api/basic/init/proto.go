package init

import (
	"flag"
	"log"
	"zhongyao/aa/crawler/practice/api/basic/config"
	__ "zhongyao/aa/crawler/practice/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitClient() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//c := pb.NewGreeterClient(conn)
	config.UserClient = __.NewUserServiceClient(conn)
}
