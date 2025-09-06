package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

var configPath = flag.String("c", "cmd/grpc/etc/config.yaml", "config file path")

func main() {
	flag.Parse()

	// 初始化应用程序
	app, err := InitializeApplication(*configPath)
	if err != nil {
		panic(err)
	}

	// 创建gRPC服务器
	srv := grpc.NewServer()

	// 注册服务
	app.RegisterServices(srv)

	// 监听端口
	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}

	// 创建上下文和信号处理
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		log.Printf("gRPC server listening on :9003")
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 等待信号
	<-sigCh
	log.Println("Shutting down gRPC server...")

	// 优雅关闭
	srv.GracefulStop()
}
