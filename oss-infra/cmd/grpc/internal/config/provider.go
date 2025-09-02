package config

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/db"
)

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig(cfg *Config) *db.DatabaseConfig {
	return &cfg.Database
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig(cfg *Config) *db.RedisConfig {
	return &cfg.Redis
}

// GetStorageConfig 获取存储配置
func GetStorageConfig(cfg *Config) *StorageConfig {
	return &cfg.Storage
}

// GetMaxFileSize 获取最大文件大小配置
func GetMaxFileSize() int64 {
	return 100 * 1024 * 1024 // 100MB
}

// GetAllowedMimeTypes 获取允许的文件类型
func GetAllowedMimeTypes() []string {
	return []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"application/pdf",
		"text/plain",
		"application/zip",
		"application/json",
	}
}

// StoragePrefix 存储前缀类型
type StoragePrefix string

// GetStoragePrefix 获取存储前缀
func GetStoragePrefix() StoragePrefix {
	return StoragePrefix("oss-files")
}

// ConfigProviderSet Config providers
var ConfigProviderSet = wire.NewSet(
	GetDatabaseConfig,
	GetRedisConfig,
	GetStorageConfig,
	GetMaxFileSize,
	GetAllowedMimeTypes,
	GetStoragePrefix,
)
