package domain

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/oss-infra/cmd/grpc/internal/config"
	"github.com/people257/poor-guy-shop/oss-infra/internal/domain/file"
)

// NewFileDomainService 创建文件领域服务
func NewFileDomainService(maxFileSize int64, allowedMimeTypes []string, storagePrefix config.StoragePrefix) *file.DomainService {
	return file.NewDomainService(maxFileSize, allowedMimeTypes, string(storagePrefix))
}

// DomainServiceProviderSet domain service provider
var DomainServiceProviderSet = wire.NewSet(
	NewFileDomainService,
)
