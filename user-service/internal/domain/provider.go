package domain

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/user-service/internal/domain/address"
	"github.com/people257/poor-guy-shop/user-service/internal/domain/auth"
	"github.com/people257/poor-guy-shop/user-service/internal/domain/user"
)

// DomainServiceProviderSet domain service provider
var DomainServiceProviderSet = wire.NewSet(
	// 用户领域服务
	user.NewDomainService,
	user.NewConverter,

	// 认证领域服务
	auth.NewDomainService,

	// 地址领域服务
	address.NewDomainService,
)
