package address

import (
	"context"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/emptypb"

	addresspb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/address"
	"github.com/people257/poor-guy-shop/user-service/internal/application/address"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/auth"
)

// AddressServer 地址gRPC服务器
type AddressServer struct {
	addresspb.UnimplementedAddressServiceServer
	addressService *address.Service
	authInfra      *auth.Auth
}

// NewAddressServer 创建地址gRPC服务器
func NewAddressServer(addressService *address.Service, authInfra *auth.Auth) *AddressServer {
	return &AddressServer{
		addressService: addressService,
		authInfra:      authInfra,
	}
}

// AddAddress 添加地址
func (s *AddressServer) AddAddress(ctx context.Context, req *addresspb.AddAddressReq) (*addresspb.AddressResp, error) {
	// 获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID()

	// 构建请求
	addReq := &address.AddAddressReq{
		UserID:        userID,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		Street:        req.Street,
		AddressLabel:  req.AddressLabel,
	}

	if req.PostalCode != "" {
		addReq.PostalCode = &req.PostalCode
	}
	if req.Longitude != "" && req.Latitude != "" {
		if lng, err := decimal.NewFromString(req.Longitude); err == nil {
			addReq.Longitude = &lng
		}
		if lat, err := decimal.NewFromString(req.Latitude); err == nil {
			addReq.Latitude = &lat
		}
	}

	// 添加地址
	resp, err := s.addressService.AddAddress(ctx, addReq)
	if err != nil {
		return nil, err
	}

	return &addresspb.AddressResp{
		Address: s.toProto(resp.Address),
	}, nil
}

// UpdateAddress 更新地址
func (s *AddressServer) UpdateAddress(ctx context.Context, req *addresspb.UpdateAddressReq) (*addresspb.AddressResp, error) {
	// 获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID()

	// 构建请求
	updateReq := &address.UpdateAddressReq{
		AddressID:     req.AddressId,
		UserID:        userID,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		Street:        req.Street,
		AddressLabel:  req.AddressLabel,
	}

	if req.PostalCode != "" {
		updateReq.PostalCode = &req.PostalCode
	}
	if req.Longitude != "" && req.Latitude != "" {
		if lng, err := decimal.NewFromString(req.Longitude); err == nil {
			updateReq.Longitude = &lng
		}
		if lat, err := decimal.NewFromString(req.Latitude); err == nil {
			updateReq.Latitude = &lat
		}
	}

	// 更新地址
	resp, err := s.addressService.UpdateAddress(ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return &addresspb.AddressResp{
		Address: s.toProto(resp.Address),
	}, nil
}

// DeleteAddress 删除地址
func (s *AddressServer) DeleteAddress(ctx context.Context, req *addresspb.DeleteAddressReq) (*emptypb.Empty, error) {
	// 获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID()

	// 删除地址
	if err := s.addressService.DeleteAddress(ctx, userID, req.AddressId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetAddresses 获取地址列表
func (s *AddressServer) GetAddresses(ctx context.Context, req *emptypb.Empty) (*addresspb.GetAddressesResp, error) {
	// 获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID()

	// 获取地址列表
	resp, err := s.addressService.GetAddresses(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为proto
	addresses := make([]*addresspb.Address, len(resp.Addresses))
	for i, addr := range resp.Addresses {
		addresses[i] = s.toProto(addr)
	}

	return &addresspb.GetAddressesResp{
		Addresses: addresses,
		Total:     int32(resp.Total),
	}, nil
}

// GetAddress 获取地址详情
func (s *AddressServer) GetAddress(ctx context.Context, req *addresspb.GetAddressReq) (*addresspb.AddressResp, error) {
	// 获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID()

	// 获取地址详情
	addr, err := s.addressService.GetAddress(ctx, userID, req.AddressId)
	if err != nil {
		return nil, err
	}

	return &addresspb.AddressResp{
		Address: s.toProto(addr),
	}, nil
}

// GetDefaultAddress 获取默认地址
func (s *AddressServer) GetDefaultAddress(ctx context.Context, req *emptypb.Empty) (*addresspb.AddressResp, error) {
	// 获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID()

	// 获取默认地址
	addr, err := s.addressService.GetDefaultAddress(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &addresspb.AddressResp{
		Address: s.toProto(addr),
	}, nil
}

// SetDefaultAddress 设置默认地址
func (s *AddressServer) SetDefaultAddress(ctx context.Context, req *addresspb.SetDefaultAddressReq) (*emptypb.Empty, error) {
	// 获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims.UserID()

	// 设置默认地址
	if err := s.addressService.SetDefaultAddress(ctx, userID, req.AddressId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// toProto 转换为proto
func (s *AddressServer) toProto(addr *address.AddressDTO) *addresspb.Address {
	protoAddr := &addresspb.Address{
		AddressId:     addr.ID,
		ReceiverName:  addr.ReceiverName,
		ReceiverPhone: addr.ReceiverPhone,
		Province:      addr.Province,
		City:          addr.City,
		District:      addr.District,
		Street:        addr.Street,
		AddressLabel:  addr.AddressLabel,
		IsDefault:     addr.IsDefault,
		FullAddress:   addr.FullAddress,
		CreatedAt:     addr.CreatedAt,
		UpdatedAt:     addr.UpdatedAt,
	}

	if addr.PostalCode != nil {
		protoAddr.PostalCode = *addr.PostalCode
	}
	if addr.Longitude != nil {
		protoAddr.Longitude = addr.Longitude.String()
	}
	if addr.Latitude != nil {
		protoAddr.Latitude = addr.Latitude.String()
	}

	return protoAddr
}
