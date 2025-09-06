package payment

import (
	"context"
	"fmt"

	"github.com/people257/poor-guy-shop/common/auth"
	pb "github.com/people257/poor-guy-shop/payment-service/gen/proto/proto/payment"
	"github.com/people257/poor-guy-shop/payment-service/internal/application/payment"
	"github.com/people257/poor-guy-shop/payment-service/internal/application/refund"
	paymentDomain "github.com/people257/poor-guy-shop/payment-service/internal/domain/payment"
	refundDomain "github.com/people257/poor-guy-shop/payment-service/internal/domain/refund"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GrpcHandler 支付gRPC处理器
type GrpcHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentService *payment.Service
	refundService  *refund.Service
}

// NewGrpcHandler 创建支付gRPC处理器
func NewGrpcHandler(
	paymentService *payment.Service,
	refundService *refund.Service,
) *GrpcHandler {
	return &GrpcHandler{
		paymentService: paymentService,
		refundService:  refundService,
	}
}

// CreatePaymentOrder 创建支付订单
func (h *GrpcHandler) CreatePaymentOrder(ctx context.Context, req *pb.CreatePaymentOrderReq) (*pb.CreatePaymentOrderResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	// 转换支付方式
	paymentMethod, err := h.convertPaymentMethod(req.PaymentMethod)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "无效的支付方式: %v", err)
	}

	// 构建请求
	serviceReq := payment.CreatePaymentOrderRequest{
		OrderID:       req.OrderId,
		Amount:        req.Amount,
		PaymentMethod: paymentMethod,
		Subject:       req.Subject,
		Description:   req.Description,
		NotifyURL:     req.NotifyUrl,
		ReturnURL:     req.ReturnUrl,
	}

	// 调用应用服务
	resp, err := h.paymentService.CreatePaymentOrder(ctx, userID, serviceReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建支付订单失败: %v", err)
	}

	// 转换响应
	return &pb.CreatePaymentOrderResp{
		PaymentOrder:  h.convertPaymentOrderToProto(resp.PaymentOrder),
		PaymentUrl:    resp.PaymentURL,
		QrCode:        resp.QRCode,
		PaymentParams: resp.PaymentParams,
	}, nil
}

// GetPaymentOrder 查询支付订单
func (h *GrpcHandler) GetPaymentOrder(ctx context.Context, req *pb.GetPaymentOrderReq) (*pb.GetPaymentOrderResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	// 调用应用服务
	paymentOrder, err := h.paymentService.GetPaymentOrder(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "支付订单不存在: %v", err)
	}

	// 检查用户权限
	if paymentOrder.UserID.String() != userID {
		return nil, status.Error(codes.PermissionDenied, "无权限访问该支付订单")
	}

	return &pb.GetPaymentOrderResp{
		PaymentOrder: h.convertPaymentOrderToProto(paymentOrder),
	}, nil
}

// HandlePaymentCallback 处理支付回调
func (h *GrpcHandler) HandlePaymentCallback(ctx context.Context, req *pb.HandlePaymentCallbackReq) (*pb.HandlePaymentCallbackResp, error) {
	// 构建请求
	serviceReq := payment.HandlePaymentCallbackRequest{
		Provider: req.Provider,
		Params:   req.Params,
	}

	// 调用应用服务
	resp, err := h.paymentService.HandlePaymentCallback(ctx, serviceReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "处理支付回调失败: %v", err)
	}

	return &pb.HandlePaymentCallbackResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// CreateRefund 申请退款
func (h *GrpcHandler) CreateRefund(ctx context.Context, req *pb.CreateRefundReq) (*pb.CreateRefundResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	// 构建请求
	serviceReq := refund.CreateRefundRequest{
		PaymentID: req.PaymentId,
		Amount:    req.Amount,
		Reason:    req.Reason,
	}

	// 调用应用服务
	resp, err := h.refundService.CreateRefund(ctx, serviceReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建退款失败: %v", err)
	}

	return &pb.CreateRefundResp{
		Refund: h.convertRefundToProto(resp.Refund),
	}, nil
}

// GetRefund 查询退款状态
func (h *GrpcHandler) GetRefund(ctx context.Context, req *pb.GetRefundReq) (*pb.GetRefundResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	// 调用应用服务
	refundEntity, err := h.refundService.GetRefund(ctx, req.RefundId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "退款记录不存在: %v", err)
	}

	return &pb.GetRefundResp{
		Refund: h.convertRefundToProto(refundEntity),
	}, nil
}

// VerifyPaymentStatus 验证支付状态
func (h *GrpcHandler) VerifyPaymentStatus(ctx context.Context, req *pb.VerifyPaymentStatusReq) (*pb.VerifyPaymentStatusResp, error) {
	// 调用应用服务
	paymentOrder, err := h.paymentService.VerifyPaymentStatus(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "支付订单不存在: %v", err)
	}

	// 转换支付状态
	paymentStatus, err := h.convertPaymentStatusToProto(paymentOrder.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "转换支付状态失败: %v", err)
	}

	return &pb.VerifyPaymentStatusResp{
		Status:    paymentStatus,
		PaymentId: paymentOrder.ID.String(),
		Amount:    paymentOrder.Amount.String(),
	}, nil
}

// convertPaymentMethod 转换支付方式
func (h *GrpcHandler) convertPaymentMethod(method pb.PaymentMethod) (paymentDomain.PaymentMethod, error) {
	switch method {
	case pb.PaymentMethod_PAYMENT_METHOD_ALIPAY:
		return paymentDomain.PaymentMethodAlipay, nil
	case pb.PaymentMethod_PAYMENT_METHOD_WECHAT:
		return paymentDomain.PaymentMethodWechat, nil
	case pb.PaymentMethod_PAYMENT_METHOD_BANK_CARD:
		return paymentDomain.PaymentMethodBankCard, nil
	case pb.PaymentMethod_PAYMENT_METHOD_BALANCE:
		return paymentDomain.PaymentMethodBalance, nil
	default:
		return "", fmt.Errorf("unsupported payment method: %v", method)
	}
}

// convertPaymentStatusToProto 转换支付状态到Proto
func (h *GrpcHandler) convertPaymentStatusToProto(status paymentDomain.PaymentStatus) (pb.PaymentStatus, error) {
	switch status {
	case paymentDomain.PaymentStatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING, nil
	case paymentDomain.PaymentStatusSuccess:
		return pb.PaymentStatus_PAYMENT_STATUS_SUCCESS, nil
	case paymentDomain.PaymentStatusFailed:
		return pb.PaymentStatus_PAYMENT_STATUS_FAILED, nil
	case paymentDomain.PaymentStatusCancelled:
		return pb.PaymentStatus_PAYMENT_STATUS_CANCELLED, nil
	case paymentDomain.PaymentStatusRefunded:
		return pb.PaymentStatus_PAYMENT_STATUS_REFUNDED, nil
	case paymentDomain.PaymentStatusPartialRefunded:
		return pb.PaymentStatus_PAYMENT_STATUS_PARTIAL_REFUNDED, nil
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED, fmt.Errorf("unknown payment status: %v", status)
	}
}

// convertPaymentOrderToProto 转换支付订单到Proto
func (h *GrpcHandler) convertPaymentOrderToProto(order *paymentDomain.PaymentOrder) *pb.PaymentOrder {
	protoOrder := &pb.PaymentOrder{
		Id:                order.ID.String(),
		OrderId:           order.OrderID.String(),
		UserId:            order.UserID.String(),
		Amount:            order.Amount.String(),
		ThirdPartyOrderId: order.ThirdPartyOrderID,
		CreatedAt:         timestamppb.New(order.CreatedAt),
	}

	// 转换支付方式
	switch order.PaymentMethod {
	case paymentDomain.PaymentMethodAlipay:
		protoOrder.PaymentMethod = pb.PaymentMethod_PAYMENT_METHOD_ALIPAY
	case paymentDomain.PaymentMethodWechat:
		protoOrder.PaymentMethod = pb.PaymentMethod_PAYMENT_METHOD_WECHAT
	case paymentDomain.PaymentMethodBankCard:
		protoOrder.PaymentMethod = pb.PaymentMethod_PAYMENT_METHOD_BANK_CARD
	case paymentDomain.PaymentMethodBalance:
		protoOrder.PaymentMethod = pb.PaymentMethod_PAYMENT_METHOD_BALANCE
	}

	// 转换支付状态
	if status, err := h.convertPaymentStatusToProto(order.Status); err == nil {
		protoOrder.Status = status
	}

	// 可选字段
	if order.PaidAt != nil {
		protoOrder.PaidAt = timestamppb.New(*order.PaidAt)
	}
	if order.ExpiredAt != nil {
		protoOrder.ExpiredAt = timestamppb.New(*order.ExpiredAt)
	}

	return protoOrder
}

// convertRefundToProto 转换退款到Proto
func (h *GrpcHandler) convertRefundToProto(refundEntity *refundDomain.Refund) *pb.Refund {
	protoRefund := &pb.Refund{
		Id:                 refundEntity.ID.String(),
		PaymentOrderId:     refundEntity.PaymentOrderID.String(),
		Amount:             refundEntity.Amount.String(),
		Reason:             refundEntity.Reason,
		ThirdPartyRefundId: refundEntity.ThirdPartyRefundID,
		CreatedAt:          timestamppb.New(refundEntity.CreatedAt),
	}

	// 转换退款状态
	switch refundEntity.Status {
	case refundDomain.RefundStatusPending:
		protoRefund.Status = pb.RefundStatus_REFUND_STATUS_PENDING
	case refundDomain.RefundStatusSuccess:
		protoRefund.Status = pb.RefundStatus_REFUND_STATUS_SUCCESS
	case refundDomain.RefundStatusFailed:
		protoRefund.Status = pb.RefundStatus_REFUND_STATUS_FAILED
	}

	// 可选字段
	if refundEntity.ProcessedAt != nil {
		protoRefund.ProcessedAt = timestamppb.New(*refundEntity.ProcessedAt)
	}

	return protoRefund
}
