package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/people257/poor-guy-shop/common/auth"
	pb "github.com/people257/poor-guy-shop/inventory-service/gen/proto/proto/inventory"
	"github.com/people257/poor-guy-shop/inventory-service/internal/application/inventory"
	"github.com/people257/poor-guy-shop/inventory-service/internal/application/reservation"
	inventoryDomain "github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
)

// Server 库存服务gRPC服务器
type Server struct {
	pb.UnimplementedInventoryServiceServer
	inventoryApp    *inventory.Service
	businessService *inventory.BusinessService
	reservationApp  *reservation.Service
}

// NewServer 创建库存服务gRPC服务器
func NewServer(inventoryApp *inventory.Service, businessService *inventory.BusinessService, reservationApp *reservation.Service) *Server {
	return &Server{
		inventoryApp:    inventoryApp,
		businessService: businessService,
		reservationApp:  reservationApp,
	}
}

// GetInventory 查询库存
func (s *Server) GetInventory(ctx context.Context, req *pb.GetInventoryReq) (*pb.GetInventoryResp, error) {
	skuID, err := uuid.Parse(req.SkuId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sku_id: %v", err)
	}

	inv, err := s.inventoryApp.GetInventory(ctx, skuID)
	if err != nil {
		if err == inventoryDomain.ErrInventoryNotFound {
			return nil, status.Errorf(codes.NotFound, "inventory not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get inventory: %v", err)
	}

	return &pb.GetInventoryResp{
		Inventory: s.inventoryToPB(inv),
	}, nil
}

// BatchGetInventory 批量查询库存
func (s *Server) BatchGetInventory(ctx context.Context, req *pb.BatchGetInventoryReq) (*pb.BatchGetInventoryResp, error) {
	if len(req.SkuIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "sku_ids cannot be empty")
	}

	skuIDs := make([]uuid.UUID, len(req.SkuIds))
	for i, skuIDStr := range req.SkuIds {
		skuID, err := uuid.Parse(skuIDStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid sku_id[%d]: %v", i, err)
		}
		skuIDs[i] = skuID
	}

	inventories, err := s.inventoryApp.BatchGetInventory(ctx, skuIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to batch get inventory: %v", err)
	}

	pbInventories := make([]*pb.Inventory, len(inventories))
	for i, inv := range inventories {
		pbInventories[i] = s.inventoryToPB(inv)
	}

	return &pb.BatchGetInventoryResp{
		Inventories: pbInventories,
	}, nil
}

// UpdateInventory 更新库存
func (s *Server) UpdateInventory(ctx context.Context, req *pb.UpdateInventoryReq) (*pb.UpdateInventoryResp, error) {
	skuID, err := uuid.Parse(req.SkuId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sku_id: %v", err)
	}

	if req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "quantity must be positive")
	}

	// 获取用户ID（如果有的话）
	var operatorID *uuid.UUID
	if userIDStr := auth.UserIDFromContext(ctx); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			operatorID = &userID
		}
	}

	changeType := s.pbToInventoryChangeType(req.Type)

	inv, err := s.inventoryApp.UpdateInventoryQuantity(ctx, skuID, changeType, req.Quantity, req.Reason, nil, operatorID)
	if err != nil {
		if err == inventoryDomain.ErrInventoryNotFound {
			return nil, status.Errorf(codes.NotFound, "inventory not found")
		}
		if err == inventoryDomain.ErrInsufficientInventory {
			return nil, status.Errorf(codes.FailedPrecondition, "insufficient inventory")
		}
		return nil, status.Errorf(codes.Internal, "failed to update inventory: %v", err)
	}

	return &pb.UpdateInventoryResp{
		Inventory: s.inventoryToPB(inv),
	}, nil
}

// ReserveInventory 预占库存
func (s *Server) ReserveInventory(ctx context.Context, req *pb.ReserveInventoryReq) (*pb.ReserveInventoryResp, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	if len(req.Items) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "items cannot be empty")
	}

	// 转换预占商品列表
	items := make([]inventoryDomain.ReserveItem, len(req.Items))
	for i, item := range req.Items {
		skuID, err := uuid.Parse(item.SkuId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid sku_id[%d]: %v", i, err)
		}
		if item.Quantity <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "quantity[%d] must be positive", i)
		}
		items[i] = inventoryDomain.ReserveItem{
			SkuID:    skuID,
			Quantity: item.Quantity,
		}
	}

	// 设置30分钟后过期
	expiresAt := time.Now().Add(30 * time.Minute)

	reservations, err := s.businessService.ReserveInventoryWithValidation(ctx, orderID, items, &expiresAt)
	if err != nil {
		if err == inventoryDomain.ErrInsufficientInventory {
			return nil, status.Errorf(codes.FailedPrecondition, "insufficient inventory")
		}
		return nil, status.Errorf(codes.Internal, "failed to reserve inventory: %v", err)
	}

	// 保存预占记录
	pbReservations := make([]*pb.InventoryReservation, len(reservations))
	for i, res := range reservations {
		// 创建预占记录
		if _, err := s.reservationApp.CreateReservation(ctx, res.SkuID, res.OrderID, res.Quantity, res.ExpiresAt); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create reservation: %v", err)
		}
		pbReservations[i] = s.reservationToPB(res)
	}

	return &pb.ReserveInventoryResp{
		Success:      true,
		Message:      "inventory reserved successfully",
		Reservations: pbReservations,
	}, nil
}

// ReleaseReservedInventory 释放预占库存
func (s *Server) ReleaseReservedInventory(ctx context.Context, req *pb.ReleaseReservedInventoryReq) (*pb.ReleaseReservedInventoryResp, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	if err := s.businessService.ReleaseInventoryWithOrderValidation(ctx, orderID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to release reservations: %v", err)
	}

	return &pb.ReleaseReservedInventoryResp{
		Success: true,
		Message: "reservations released successfully",
	}, nil
}

// ConfirmInventoryDeduction 确认扣减库存
func (s *Server) ConfirmInventoryDeduction(ctx context.Context, req *pb.ConfirmInventoryDeductionReq) (*pb.ConfirmInventoryDeductionResp, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_id: %v", err)
	}

	if err := s.businessService.ConfirmInventoryWithOrderValidation(ctx, orderID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to confirm reservations: %v", err)
	}

	return &pb.ConfirmInventoryDeductionResp{
		Success: true,
		Message: "inventory deduction confirmed successfully",
	}, nil
}

// GetInventoryLogs 库存变动日志
func (s *Server) GetInventoryLogs(ctx context.Context, req *pb.GetInventoryLogsReq) (*pb.GetInventoryLogsResp, error) {
	skuID, err := uuid.Parse(req.SkuId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sku_id: %v", err)
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	logs, total, err := s.inventoryApp.GetInventoryLogs(ctx, skuID, int(page), int(pageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get inventory logs: %v", err)
	}

	pbLogs := make([]*pb.InventoryLog, len(logs))
	for i, log := range logs {
		pbLogs[i] = s.inventoryLogToPB(log)
	}

	return &pb.GetInventoryLogsResp{
		Logs:     pbLogs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CheckInventoryAvailability 检查库存充足性（内部RPC）
func (s *Server) CheckInventoryAvailability(ctx context.Context, req *pb.CheckInventoryAvailabilityReq) (*pb.CheckInventoryAvailabilityResp, error) {
	if len(req.Items) == 0 {
		return &pb.CheckInventoryAvailabilityResp{
			Available:        true,
			InsufficientSkus: nil,
		}, nil
	}

	// 转换检查商品列表
	items := make([]inventoryDomain.ReserveItem, len(req.Items))
	for i, item := range req.Items {
		skuID, err := uuid.Parse(item.SkuId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid sku_id[%d]: %v", i, err)
		}
		if item.Quantity <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "quantity[%d] must be positive", i)
		}
		items[i] = inventoryDomain.ReserveItem{
			SkuID:    skuID,
			Quantity: item.Quantity,
		}
	}

	available, insufficientSkus, err := s.inventoryApp.CheckInventoryAvailability(ctx, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check inventory availability: %v", err)
	}

	insufficientSkuStrings := make([]string, len(insufficientSkus))
	for i, skuID := range insufficientSkus {
		insufficientSkuStrings[i] = skuID.String()
	}

	return &pb.CheckInventoryAvailabilityResp{
		Available:        available,
		InsufficientSkus: insufficientSkuStrings,
	}, nil
}

// 辅助方法：将领域对象转换为protobuf对象
func (s *Server) inventoryToPB(inv *inventoryDomain.Inventory) *pb.Inventory {
	return &pb.Inventory{
		SkuId:             inv.SkuID.String(),
		AvailableQuantity: inv.AvailableQuantity,
		ReservedQuantity:  inv.ReservedQuantity,
		TotalQuantity:     inv.TotalQuantity,
		AlertQuantity:     inv.AlertQuantity,
		UpdatedAt:         timestamppb.New(inv.UpdatedAt),
	}
}

func (s *Server) inventoryLogToPB(log *inventoryDomain.InventoryLog) *pb.InventoryLog {
	pbLog := &pb.InventoryLog{
		Id:             log.ID.String(),
		SkuId:          log.SkuID.String(),
		Type:           s.inventoryChangeTypeToPB(log.Type),
		Quantity:       log.Quantity,
		BeforeQuantity: log.BeforeQuantity,
		AfterQuantity:  log.AfterQuantity,
		Reason:         log.Reason,
		CreatedAt:      timestamppb.New(log.CreatedAt),
	}

	if log.OrderID != nil {
		pbLog.OrderId = log.OrderID.String()
	}

	return pbLog
}

func (s *Server) reservationToPB(res *inventoryDomain.InventoryReservation) *pb.InventoryReservation {
	pbRes := &pb.InventoryReservation{
		Id:        res.ID.String(),
		SkuId:     res.SkuID.String(),
		OrderId:   res.OrderID.String(),
		Quantity:  res.Quantity,
		Status:    string(res.Status),
		CreatedAt: timestamppb.New(res.CreatedAt),
	}

	if res.ExpiresAt != nil {
		pbRes.ExpiresAt = timestamppb.New(*res.ExpiresAt)
	}

	return pbRes
}

func (s *Server) pbToInventoryChangeType(pbType pb.InventoryChangeType) inventoryDomain.InventoryChangeType {
	switch pbType {
	case pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_IN:
		return inventoryDomain.InventoryChangeTypeIn
	case pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_OUT:
		return inventoryDomain.InventoryChangeTypeOut
	case pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_RESERVE:
		return inventoryDomain.InventoryChangeTypeReserve
	case pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_RELEASE:
		return inventoryDomain.InventoryChangeTypeRelease
	case pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_ADJUST:
		return inventoryDomain.InventoryChangeTypeAdjust
	default:
		return inventoryDomain.InventoryChangeTypeAdjust
	}
}

func (s *Server) inventoryChangeTypeToPB(changeType inventoryDomain.InventoryChangeType) pb.InventoryChangeType {
	switch changeType {
	case inventoryDomain.InventoryChangeTypeIn:
		return pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_IN
	case inventoryDomain.InventoryChangeTypeOut:
		return pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_OUT
	case inventoryDomain.InventoryChangeTypeReserve:
		return pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_RESERVE
	case inventoryDomain.InventoryChangeTypeRelease:
		return pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_RELEASE
	case inventoryDomain.InventoryChangeTypeAdjust:
		return pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_ADJUST
	default:
		return pb.InventoryChangeType_INVENTORY_CHANGE_TYPE_UNSPECIFIED
	}
}
