package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 初始化应用程序
	app, cleanup, err := InitializeApplication(ctx, "etc/config.yaml")
	if err != nil {
		log.Fatalf("初始化应用程序失败: %v", err)
	}
	defer cleanup()

	// 监听系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动应用程序
	go func() {
		log.Println("启动订单服务...")
		if err := app.Run(ctx); err != nil {
			log.Fatalf("启动服务失败: %v", err)
		}
	}()

	// 等待关闭信号
	select {
	case sig := <-sigCh:
		log.Printf("接收到信号: %v, 开始优雅关闭...", sig)
		cancel()
	case <-ctx.Done():
		log.Println("上下文取消，开始关闭...")
	}

	log.Println("订单服务已关闭")
}
