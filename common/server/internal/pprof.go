package internal

import (
	"fmt"
	"github.com/people257/poor-guy-shop/common/server/config"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
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
