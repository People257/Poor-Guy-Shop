package main

import (
	"context"

	"github.com/people257/poor-guy-shop/common/server"
	"github.com/people257/poor-guy-shop/product-service/api/brand"
	"github.com/people257/poor-guy-shop/product-service/api/category"
	"github.com/people257/poor-guy-shop/product-service/api/product"
)

// Application 应用程序
type Application struct {
	server *server.Server
}

// NewApplication 创建应用程序
func NewApplication(
	srv *server.Server,
	categoryServer *category.CategoryServer,
	brandServer *brand.BrandServer,
	productServer *product.ProductServer,
) *Application {
	return &Application{
		server: srv,
	}
}

// Run 运行应用程序
func (a *Application) Run(ctx context.Context) error {
	return a.server.Run(ctx)
}
