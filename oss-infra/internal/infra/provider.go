package infra

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/oss-infra/internal/infra/repository"
	"github.com/people257/poor-guy-shop/oss-infra/internal/infra/storage"
)

// InfraProviderSet Infrastructure providers
var InfraProviderSet = wire.NewSet(
	repository.NewFileRepository,
	storage.NewStorageRepository,
)
