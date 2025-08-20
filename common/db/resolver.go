package db

import (
	"fmt"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// ReplicaConfig 从库配置
type ReplicaConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Port     uint16 `mapstructure:"port"`
}

// setupDBResolver 配置数据库读写分离
func setupDBResolver(db *gorm.DB, cfg *DatabaseConfig) error {
	// 配置从库连接
	var replicaDialectors []gorm.Dialector
	for _, replica := range cfg.Replicas {
		replicaDSN := fmt.Sprintf(postgresTcpDSN,
			replica.Host,
			replica.User,
			replica.Password,
			replica.Database,
			replica.Port,
		)
		replicaDialectors = append(replicaDialectors, postgres.Open(replicaDSN))
	}

	// 注册 DBResolver
	resolverConfig := dbresolver.Config{
		Replicas:          replicaDialectors,
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: true,
	}

	// 配置连接池和读写分离
	return db.Use(
		dbresolver.Register(resolverConfig).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)
}
