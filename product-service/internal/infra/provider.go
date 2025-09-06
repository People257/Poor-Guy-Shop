package infra

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/product-service/internal/infra/repository"
)

// ProviderSet 基础设施层依赖注入提供者集合
var ProviderSet = wire.NewSet(
	repository.NewCategoryRepository,
	repository.NewBrandRepository,
	repository.NewProductRepository,
	repository.NewProductSKURepository,
)
