package storage

import (
	"fmt"

	"github.com/people257/poor-guy-shop/oss-infra/internal/domain/file"
)

// StorageType 存储类型
type StorageType string

const (
	StorageTypeS3  StorageType = "s3"
	StorageTypeOSS StorageType = "oss"
	StorageTypeCOS StorageType = "cos"
)

// Config 存储配置
type Config struct {
	Type StorageType `yaml:"type" json:"type"` // 存储类型

	// 云存储配置
	Region          string `yaml:"region" json:"region"`
	Bucket          string `yaml:"bucket" json:"bucket"`
	AccessKeyID     string `yaml:"access_key_id" json:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key" json:"secret_access_key"`
	Endpoint        string `yaml:"endpoint" json:"endpoint"`

	// 高级配置
	MaxFileSize      int64    `yaml:"max_file_size" json:"max_file_size"`           // 最大文件大小
	AllowedMimeTypes []string `yaml:"allowed_mime_types" json:"allowed_mime_types"` // 允许的MIME类型
	StoragePrefix    string   `yaml:"storage_prefix" json:"storage_prefix"`         // 存储前缀
}

// NewStorageRepository 根据配置创建存储仓储
func NewStorageRepository(cfg *config.StorageConfig) (file.StorageRepository, error) {
	switch cfg.Provider {
	case "aliyun":
		return NewAliyunOSSStorage(&cfg.Aliyun)
	case "qiniu":
		return NewQiniuKodoStorage(&cfg.Qiniu)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", cfg.Provider)
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	switch c.Type {
	case StorageTypeS3:
		return c.validateS3Config()
	case StorageTypeOSS:
		return c.validateOSSConfig()
	case StorageTypeCOS:
		return c.validateCOSConfig()
	default:
		return fmt.Errorf("无效的存储类型: %s", c.Type)
	}
}

// validateS3Config 验证S3配置
func (c *Config) validateS3Config() error {
	if c.Bucket == "" {
		return fmt.Errorf("S3 bucket不能为空")
	}
	if c.AccessKeyID == "" {
		return fmt.Errorf("S3 AccessKeyID不能为空")
	}
	if c.SecretAccessKey == "" {
		return fmt.Errorf("S3 SecretAccessKey不能为空")
	}
	if c.Region == "" {
		c.Region = "us-east-1" // 默认区域
	}
	return nil
}

// validateOSSConfig 验证阿里云OSS配置
func (c *Config) validateOSSConfig() error {
	if c.Bucket == "" {
		return fmt.Errorf("OSS bucket不能为空")
	}
	if c.AccessKeyID == "" {
		return fmt.Errorf("OSS AccessKeyID不能为空")
	}
	if c.SecretAccessKey == "" {
		return fmt.Errorf("OSS SecretAccessKey不能为空")
	}
	if c.Endpoint == "" {
		return fmt.Errorf("OSS Endpoint不能为空")
	}
	return nil
}

// validateCOSConfig 验证腾讯云COS配置
func (c *Config) validateCOSConfig() error {
	if c.Bucket == "" {
		return fmt.Errorf("COS bucket不能为空")
	}
	if c.AccessKeyID == "" {
		return fmt.Errorf("COS AccessKeyID不能为空")
	}
	if c.SecretAccessKey == "" {
		return fmt.Errorf("COS SecretAccessKey不能为空")
	}
	if c.Region == "" {
		return fmt.Errorf("COS Region不能为空")
	}
	return nil
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Type:             StorageTypeS3,
		MaxFileSize:      100 * 1024 * 1024, // 100MB
		AllowedMimeTypes: file.GetDefaultAllowedMimeTypes(),
		StoragePrefix:    "oss",
	}
}
