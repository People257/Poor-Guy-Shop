package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

type Config struct {
	GrpcServerConfig config.GrpcServerConfig `mapstructure:",squash"`
	Database         db.DatabaseConfig       `mapstructure:"database"`
	Redis            db.RedisConfig          `mapstructure:"redis"`
	Storage          StorageConfig           `mapstructure:"storage"`
}

type StorageConfig struct {
	Provider string              `mapstructure:"provider"` // aliyun, qiniu
	Aliyun   AliyunStorageConfig `mapstructure:"aliyun"`
	Qiniu    QiniuStorageConfig  `mapstructure:"qiniu"`
}

type AliyunStorageConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	Endpoint        string `mapstructure:"endpoint"`
}

type QiniuStorageConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"` // 自定义域名或默认域名
	Zone      string `mapstructure:"zone"`   // 存储区域
}

func MustLoad(path string) *Config {
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		panic(err)
	}
	var cfg Config
	unmarshalConfig := koanf.UnmarshalConf{
		Tag:       "mapstructure",
		FlatPaths: false,
	}
	if err := k.UnmarshalWithConf("", &cfg, unmarshalConfig); err != nil {
		panic(err)
	}
	return &cfg
}

func GetGrpcServerConfig(cfg *Config) *config.GrpcServerConfig {
	if cfg == nil {
		panic("grpc server config is nil")
	}
	return &cfg.GrpcServerConfig
}
