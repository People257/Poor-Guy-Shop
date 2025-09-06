package domain

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/product-service/internal/domain/brand"
	"github.com/people257/poor-guy-shop/product-service/internal/domain/category"
	"github.com/people257/poor-guy-shop/product-service/internal/domain/product"
)

// ProviderSet 领域服务依赖注入提供者集合
var ProviderSet = wire.NewSet(
	category.NewDomainService,
	brand.NewDomainService,
	product.NewDomainService,
)
