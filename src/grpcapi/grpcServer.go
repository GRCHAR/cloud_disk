package grpc

import (
	"cloud_disk/src/grpc/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type HelloServiceImpl struct {
}

func (p *HelloServiceImpl) Hello(ctx context.Context, args *proto.String) (*proto.String, error) {
	reply := &proto.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func init() {
	grpcServer := grpc.NewServer()
	proto.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))
	listener, err := net.Listen("tcp", ":6869")
	if err != nil {
		log.Println(err)
		return
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("grpc服务启动！！！")
}

func StartGrpcServer() {

}
