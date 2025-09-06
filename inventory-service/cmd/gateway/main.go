package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/people257/poor-guy-shop/inventory-service/gen/proto/proto/inventory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9003", "gRPC server endpoint")
		httpPort           = flag.Int("http-port", 8003, "HTTP server port")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建gRPC-Gateway多路复用器
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 注册库存服务
	err := pb.RegisterInventoryServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register inventory service handler: %v", err)
	}

	// 启动HTTP服务器
	httpAddr := fmt.Sprintf(":%d", *httpPort)
	log.Printf("Starting HTTP gateway server on %s", httpAddr)
	log.Printf("Proxying to gRPC server at %s", *grpcServerEndpoint)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatalf("Failed to serve HTTP gateway: %v", err)
	}
}
