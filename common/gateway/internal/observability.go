package internal

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"poor-guy-shop/common/gateway/config"
)

func NewObservabilityHttpServer(cfg *config.ObservabilityConfig) (*http.Server, func()) {
	httpSrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port)}
	m := http.NewServeMux()

	if cfg.Pprof.Enable {
		m.HandleFunc("GET /debug/pprof/", pprof.Index)
		m.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		m.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		m.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		m.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	httpSrv.Handler = m
	cleanUp := func() {
		if err := httpSrv.Close(); err != nil {
			zap.L().Error("failed to close http server", zap.Error(err))
		}
	}
	return httpSrv, cleanUp
}
