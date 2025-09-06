package infra

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/payment-service/internal/infra/payment"
	"github.com/people257/poor-guy-shop/payment-service/internal/infra/repository"
)

// ProviderSet 基础设施层提供者集合
var ProviderSet = wire.NewSet(
	// 仓储层
	repository.NewPaymentRepository,
	repository.NewRefundRepository,

	// 支付服务
	payment.NewAlipayClient,
)
