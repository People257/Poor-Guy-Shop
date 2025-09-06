package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbcart "github.com/people257/poor-guy-shop/order-service/gen/proto/order/cart"
	pborder "github.com/people257/poor-guy-shop/order-service/gen/proto/order/order"
)

func main() {
	var (
		grpcAddr    = flag.String("grpc-addr", ":9002", "gRPC server address")
		gatewayAddr = flag.String("gateway-addr", ":8002", "Gateway server address")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建gRPC连接
	conn, err := grpc.DialContext(ctx, *grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	// 创建gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// 注册订单服务
	if err := pborder.RegisterOrderServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalf("Failed to register order service handler: %v", err)
	}

	// 注册购物车服务
	if err := pbcart.RegisterCartServiceHandler(ctx, mux, conn); err != nil {
		log.Fatalf("Failed to register cart service handler: %v", err)
	}

	log.Printf("Starting gateway server on %s", *gatewayAddr)
	if err := http.ListenAndServe(*gatewayAddr, mux); err != nil {
		log.Fatalf("Failed to start gateway server: %v", err)
	}
}

