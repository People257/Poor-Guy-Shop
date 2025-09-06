package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var configPath = flag.String("f", "etc/config.yaml", "config file path")

func main() {
	flag.Parse()

	fmt.Println("启动用户服务...")
	fmt.Printf("配置文件路径: %s\n", *configPath)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)

	fmt.Println("正在初始化应用...")
	app, cleanUp := InitializeApplication(ctx, *configPath)
	fmt.Println("应用初始化完成")

	errGroup.Go(func() error {
		fmt.Println("正在启动服务...")
		return app.Run(ctx)
	})
	errGroup.Go(func() error {
		<-ctx.Done()
		cleanUp()
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		zap.L().Error("exit with error", zap.Error(err))
	}
}
