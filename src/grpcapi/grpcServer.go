package grpc

import (
	"cloud_disk/src/grpcapi/proto"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type HelloServiceImpl struct {
	proto.UnimplementedHelloServiceServer
}

var serverClientStartChannel = make(chan bool, 1)

func (*HelloServiceImpl) Hello(ctx context.Context, args *proto.String) (*proto.String, error) {
	reply := &proto.String{Value: args.Value}
	log.Println("Hello message :", args.Value)
	return reply, nil
}

func (*HelloServiceImpl) HelloMessage(ctx context.Context, args *proto.Message) (*proto.Message, error) {
	reply := &proto.Message{Value: args.Value}
	log.Println("Hello message :", args.Value)
	time.Sleep(5 * time.Second)
	return reply, nil
}

func (*HelloServiceImpl) HelloTwo(ctx context.Context, args *proto.String) (*proto.String, error) {
	reply := &proto.String{Value: args.Value}
	return reply, nil
}

func (*HelloServiceImpl) HelloServerStream(args *proto.Count, stream proto.HelloService_HelloServerStreamServer) (err error) {
	timeCount := int64(0)
	for {
		select {
		case <-time.After(1 * time.Second):
			timeCount++
			reply := &proto.Count{Value: args.Value + timeCount}
			err := stream.Send(reply)
			if err != nil {
				log.Println("send error:", err)
				return err
			}
		}
	}
	return nil
}

func (*HelloServiceImpl) HelloClientStream(stream proto.HelloService_HelloClientStreamServer) (err error) {
	return nil
}

func (*HelloServiceImpl) HelloDoubleStream(stream proto.HelloService_HelloDoubleStreamServer) (err error) {
	countChannel := make(chan int64, 10000)
	errChannel := make(chan error, 1)
	go func() {
		for {
			count, err := stream.Recv()

			if err != nil {
				log.Println(err)
				if err != io.EOF {
					errChannel <- err
				}
				return
			}
			log.Println("recv count:", count.Value)
			countChannel <- count.Value + 100
		}

	}()
	for {
		select {
		case receivedCount := <-countChannel:
			log.Println("send count:", receivedCount)
			reply := &proto.Count{Value: receivedCount}
			err := stream.Send(reply)
			if err != nil {
				return err
			}
		case err := <-errChannel:
			return err
		}
	}

	return nil
}

func (*HelloServiceImpl) mustEmbedUnimplementedHelloServiceServer() {}

func init() {

}

func StartGrpcServer() {
	grpcServer := grpc.NewServer()
	proto.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))
	listener, err := net.Listen("tcp", "0.0.0.0:6880")

	if err != nil {
		log.Println("grpc接口监听失败:", err)
		return
	}
	log.Println("grpc服务开始监听")
	serverClientStartChannel <- true
	go func() {

		err = grpcServer.Serve(listener)
		if err != nil {
			log.Println("grpc服务启动失败：", err)
			return
		}
		log.Println("grpc服务启动！！！")

	}()

}
