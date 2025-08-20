package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var configPath = flag.String("f", "etc/config.yaml", "config file path")

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)

	app, cleanUp := InitializeApplication(ctx, *configPath)

	errGroup.Go(func() error {
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
