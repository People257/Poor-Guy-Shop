package client

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/people257/poor-guy-shop/order-service/cmd/grpc/config"
)

// PaymentServiceClient 支付服务客户端
type PaymentServiceClient struct {
	config *config.ServiceConfig
	conn   *grpc.ClientConn
}

// NewPaymentServiceClient 创建支付服务客户端
func NewPaymentServiceClient(cfg *config.ServiceConfig) (*PaymentServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to payment service at %s: %v. Using mock implementation.", addr, err)
		// 连接失败时使用mock模式
		return &PaymentServiceClient{
			config: cfg,
			conn:   nil,
		}, nil
	}

	return &PaymentServiceClient{
		config: cfg,
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (c *PaymentServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreatePayment 创建支付订单
func (c *PaymentServiceClient) CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	if c.conn == nil {
		// Mock模式：直接返回成功
		log.Printf("Payment client in mock mode: creating payment for order %s", req.OrderID)
		return &PaymentResponse{
			Success:    true,
			PaymentURL: fmt.Sprintf("http://mock-payment.com/pay/%s", req.OrderID),
			QRCode:     "mock_qr_code",
			PaymentParams: map[string]string{
				"order_id": req.OrderID,
				"amount":   req.Amount,
			},
		}, nil
	}

	// TODO: 实现真实的gRPC调用
	// 由于proto定义依赖问题，暂时使用mock
	log.Printf("Payment service connected but using mock implementation for order %s", req.OrderID)
	return &PaymentResponse{
		Success:    true,
		PaymentURL: fmt.Sprintf("http://connected-payment.com/pay/%s", req.OrderID),
		QRCode:     "connected_qr_code",
		PaymentParams: map[string]string{
			"order_id": req.OrderID,
			"amount":   req.Amount,
		},
	}, nil
}

// VerifyPayment 验证支付状态
func (c *PaymentServiceClient) VerifyPayment(ctx context.Context, orderID string) (bool, error) {
	if c.conn == nil {
		log.Printf("Payment client in mock mode: payment verified for order %s", orderID)
		return true, nil
	}

	// TODO: 实现真实的gRPC调用
	log.Printf("Payment service connected but using mock verification for order %s", orderID)
	return true, nil
}
