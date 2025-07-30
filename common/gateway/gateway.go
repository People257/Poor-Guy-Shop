package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net/http"
	"poor-guy-shop/common/gateway/config"
)

type Gateway struct {
	Config            *config.GatewayConfig
	Echo              *echo.Echo
	metricsHttpServer *http.Server
	conn              *grpc.ClientConn
	gwmux             *runtime.ServeMux
}

func newGateway(
	echo *echo.Echo,
	config *config.GatewayConfig,
	metricsHttpServer *http.Server,
	conn *grpc.ClientConn,
	gwmux *runtime.ServeMux,

	_ trace.TracerProvider,
	_ metric.MeterProvider,
) *Gateway {
	return &Gateway{Echo: echo, metricsHttpServer: metricsHttpServer, Config: config, conn: conn, gwmux: gwmux}
}

func (g *Gateway) RegisterHandler(fn func(gwmux *runtime.ServeMux, conn *grpc.ClientConn) error) error {
	return fn(g.gwmux, g.conn)
}

func (g *Gateway) startGateway() error {
	serverCfg := g.Config.Server
	zap.L().Info("starting gateway server", zap.Uint16("port", serverCfg.Port))
	if err := g.Echo.Start(fmt.Sprintf(":%d", serverCfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (g *Gateway) startObservabilityHttp() error {
	isMetricsEnable := g.Config.Observability.Metrics.Enable
	isPprofEnable := g.Config.Observability.Pprof.Enable

	if !isMetricsEnable && !isPprofEnable {
		return nil
	}

	zap.L().Info("starting observability http server", zap.Uint16("port", g.Config.Observability.Port),
		zap.Bool("metrics", isMetricsEnable),
		zap.Bool("pprof", isPprofEnable))
	if err := g.metricsHttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (g *Gateway) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return g.startGateway()
	})

	errGroup.Go(func() error {
		return g.startObservabilityHttp()
	})

	return errGroup.Wait()
}
