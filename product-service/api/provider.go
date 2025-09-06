package api

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/product-service/api/brand"
	"github.com/people257/poor-guy-shop/product-service/api/category"
	"github.com/people257/poor-guy-shop/product-service/api/product"
)

// ProviderSet API层依赖注入提供者集合
var ProviderSet = wire.NewSet(
	category.NewCategoryServer,
	brand.NewBrandServer,
	product.NewProductServer,
)
