package internal

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"poor-guy-shop/common/gateway/config"
	"poor-guy-shop/common/resolver"
	"time"
)

var DefaultTimeout = 5 * time.Second

func NewGatewayMux() *runtime.ServeMux {
	gwmux := createGatewayMux()
	return gwmux
}

func NewGrpcClient(registryCfg *config.RegistryConfig) (*grpc.ClientConn, func()) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s/%s", registryCfg.Address, registryCfg.Service),
		grpc.WithResolvers(&resolver.Builder{}),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	cleanUp := func() {
		if err := conn.Close(); err != nil {
			zap.L().Error("failed to close grpc conn", zap.Error(err))
		}
	}
	return conn, cleanUp
}

func NewEcho(gwmux *runtime.ServeMux, logger *zap.Logger, serverCfg *config.ServerConfig) (*echo.Echo, func()) {
	if serverCfg.Env == config.EnvProd {
		runtime.DefaultContextTimeout = DefaultTimeout
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.JSONSerializer = &EchoSonicJSONSerializer{}

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll:   true,
		DisablePrintStack: true,
		LogErrorFunc: func(c echo.Context, err error, _ []byte) error {
			zap.L().Error("recovered panic", zap.Error(err))
			return nil
		},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetPath(c.Request().URL.Path)
			return next(c)
		}
	})
	e.Use(otelecho.Middleware(serverCfg.Name))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogUserAgent: true,
		LogURIPath:   true,

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			traceID := trace.SpanContextFromContext(c.Request().Context()).TraceID().String()
			logger.Info("request",
				zap.String("trace_id", traceID),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.String("method", v.Method),
				zap.Duration("latency_human", v.Latency),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("user_agent", v.UserAgent),
				zap.Int64("latency_ms", v.Latency.Milliseconds()),
				zap.String("uri_path", v.URIPath),
			)

			return nil
		},
	}))

	e.Any("/*", echo.WrapHandler(gwmux))

	cleanUp := func() {
		if err := e.Shutdown(context.Background()); err != nil {
			zap.L().Error("failed to shutdown echo", zap.Error(err))
		}
	}

	return e, cleanUp
}

func createGatewayMux() *runtime.ServeMux {
	gwmux := runtime.NewServeMux(
		runtime.WithErrorHandler(ErrorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, NewRespMarshaler()),
		runtime.WithForwardResponseRewriter(ResponseRewriter),
	)
	return gwmux
}
