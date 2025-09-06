package application

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/payment-service/internal/application/payment"
	"github.com/people257/poor-guy-shop/payment-service/internal/application/refund"
	paymentDomain "github.com/people257/poor-guy-shop/payment-service/internal/domain/payment"
	refundDomain "github.com/people257/poor-guy-shop/payment-service/internal/domain/refund"
)

// ProviderSet 应用层提供者集合
var ProviderSet = wire.NewSet(
	// 领域服务
	paymentDomain.NewDomainService,
	refundDomain.NewDomainService,
	
	// 应用服务
	payment.NewService,
	refund.NewService,
)

