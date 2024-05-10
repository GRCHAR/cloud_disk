package main

import (
	grpc "cloud_disk/src/grpcapi"
	"log"
	"sync"
)

func main() {
	log.Println("cloud disk start!")
	//router.InitRouter()
	grpc.StartGrpcServer()
	//grpc.StartGrpcClient()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
