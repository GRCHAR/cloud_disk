package grpc

import (
	"cloud_disk/src/grpcapi/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

func init() {

}

func StartGrpcClient() {
	if start := <-serverClientStartChannel; start {
		log.Println("start grpc client")
	}
	conn, err := grpc.Dial("localhost:6880", grpc.WithInsecure())
	if err != nil {
		log.Fatal("grpc dial error:", err)
	}
	//defer func(conn *grpc.ClientConn) {
	//	err := conn.Close()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}(conn)
	client := proto.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &proto.String{Value: "hello"})
	if err != nil {
		log.Fatal("grpc hello error:", err)
	}
	log.Println("server reply:", reply)

	stream, err := client.HelloServerStream(context.Background(), &proto.Count{
		Value: 0,
	})
	if err != nil {
		log.Println("HelloServerStream error:", err)
		return
	}
	go func() {
		for {
			recv, err := stream.Recv()
			if err != nil {
				return
			}
			log.Println("server receive:", recv.Value)
		}
	}()
}
