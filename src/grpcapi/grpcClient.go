package grpc

import (
	"cloud_disk/src/grpcapi/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

func init() {
	conn, err := grpc.Dial("localhost:6869", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := proto.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &proto.String{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server reply:", reply)
}

func StartGrpcClient() {

}
