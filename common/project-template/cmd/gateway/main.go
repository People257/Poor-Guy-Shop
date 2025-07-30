package main

import (
	_ "cnb.cool/cymirror/ces-services/common/resolver"
	"context"
	"flag"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
)

var configPath = flag.String("f", "etc/config.yaml", "config file path")

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)

	gateway, cleanUp := InitializeApplication(ctx, *configPath)

	errGroup.Go(func() error {
		return gateway.Run(ctx)
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
