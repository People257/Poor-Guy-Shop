package db

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
)

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	TLS      bool   `mapstructure:"tls"`
	Cluster  bool   `mapstructure:"cluster"`
}

func NewRedis(cfg *RedisConfig) redis.UniversalClient {
	var tlsConfig *tls.Config
	if cfg.TLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	var client redis.UniversalClient

	if cfg.Cluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:     []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
			Username:  cfg.User,
			Password:  cfg.Password,
			TLSConfig: tlsConfig,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Username:  cfg.User,
			Password:  cfg.Password,
			TLSConfig: tlsConfig,
		})
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	err := redisotel.InstrumentTracing(client, redisotel.WithTracerProvider(otel.GetTracerProvider()))
	if err != nil {
		panic(err)
	}

	err = redisotel.InstrumentMetrics(client, redisotel.WithMeterProvider(otel.GetMeterProvider()))
	if err != nil {
		panic(err)
	}

	return client
}
