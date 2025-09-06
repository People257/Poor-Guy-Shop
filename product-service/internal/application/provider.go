package application

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/product-service/internal/application/brand"
	"github.com/people257/poor-guy-shop/product-service/internal/application/category"
	"github.com/people257/poor-guy-shop/product-service/internal/application/product"
)

// ProviderSet 应用层依赖注入提供者集合
var ProviderSet = wire.NewSet(
	category.NewService,
	brand.NewService,
	product.NewService,
)
