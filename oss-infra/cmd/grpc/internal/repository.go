package internal

import (
	"fmt"

	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/oss-infra/gen/gen/query"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"moul.io/zapgorm2"
)

const postgresTcpDSN = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai"

func NewDB(cfg *db.DatabaseConfig) *db.DB[*query.Query] {
	logger := zapgorm2.New(zap.L())
	logger.IgnoreRecordNotFoundError = true

	dsn := fmt.Sprintf(postgresTcpDSN, cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger,
		TranslateError:         true,
	})
	if err != nil {
		panic(err)
	}
	err = database.Use(tracing.NewPlugin(tracing.WithTracerProvider(otel.GetTracerProvider())))
	if err != nil {
		panic(err)
	}

	return db.New(query.Use(database))
}

func NewRedisClient(cfg *db.RedisConfig) redis.UniversalClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.User,
		Password: cfg.Password,
	})
	err := redisotel.InstrumentTracing(client, redisotel.WithTracerProvider(otel.GetTracerProvider()))
	if err != nil {
		panic(err)
	}

	return client
}
