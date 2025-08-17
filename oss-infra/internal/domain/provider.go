package domain

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/oss-infra/internal/domain/file"
)

// DomainServiceProviderSet domain service provider
var DomainServiceProviderSet = wire.NewSet(
	file.NewDomainService,
)
