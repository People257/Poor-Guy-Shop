package db

import (
	"context"
	"database/sql"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"moul.io/zapgorm2"
)

const postgresTcpDSN = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai"

var key keyType = struct{}{}

type keyType struct{}

type Transactional[T any] interface {
	Transaction(fn func(tx T) error, options ...*sql.TxOptions) error
}

type DB[T Transactional[T]] struct {
	query T
}

var instance any

func newDB[T Transactional[T]](query T) *DB[T] {
	database := &DB[T]{query: query}
	instance = database
	return database
}

func (d *DB[T]) Get(ctx context.Context) T {
	db, ok := ctx.Value(key).(T)
	if !ok {
		return d.query
	}
	return db
}

func Transaction[T Transactional[T]](ctx context.Context, fn func(ctx context.Context) error, options ...*sql.TxOptions) error {
	db := instance.(*DB[T])
	return db.query.Transaction(func(tx T) error {
		ctx = context.WithValue(ctx, key, tx)
		return fn(ctx)
	}, options...)
}

// DatabaseConfig 主从数据库配置
type DatabaseConfig struct {
	// 主库配置
	Master struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
		Port     uint16 `mapstructure:"port"`
	} `mapstructure:"master"`
	// 从库配置
	Replicas []ReplicaConfig `mapstructure:"replicas"`
}

func NewDB[T Transactional[T]](cfg *DatabaseConfig,
	newQueryFunc func(db *gorm.DB, opts ...gen.DOOption) T) *DB[T] {
	logger := zapgorm2.New(zap.L())
	logger.IgnoreRecordNotFoundError = true

	// 配置主库 DSN
	masterDSN := fmt.Sprintf(postgresTcpDSN,
		cfg.Master.Host,
		cfg.Master.User,
		cfg.Master.Password,
		cfg.Master.Database,
		cfg.Master.Port,
	)

	// 初始化主库连接
	database, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger,
		TranslateError:         true,
	})
	if err != nil {
		panic(err)
	}

	// 配置读写分离
	if err = setupDBResolver(database, cfg); err != nil {
		panic(err)
	}

	// 配置 OpenTelemetry
	if err = database.Use(tracing.NewPlugin(
		tracing.WithTracerProvider(otel.GetTracerProvider()),
		tracing.WithoutServerAddress(),
	)); err != nil {
		panic(err)
	}

	return newDB[T](newQueryFunc(database))
}
