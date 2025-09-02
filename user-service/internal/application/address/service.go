package address

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	"github.com/people257/poor-guy-shop/user-service/internal/domain/address"
)

// Service 地址应用服务
type Service struct {
	addressDomain *address.DomainService
	addressRepo   address.Repository
}

// NewService 创建地址应用服务
func NewService(addressDomain *address.DomainService, addressRepo address.Repository) *Service {
	return &Service{
		addressDomain: addressDomain,
		addressRepo:   addressRepo,
	}
}

// AddAddressReq 添加地址请求
type AddAddressReq struct {
	UserID        string
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	Street        string
	PostalCode    *string
	AddressLabel  string
	Longitude     *decimal.Decimal
	Latitude      *decimal.Decimal
}

// AddAddressResp 添加地址响应
type AddAddressResp struct {
	Address *AddressDTO
}

// UpdateAddressReq 更新地址请求
type UpdateAddressReq struct {
	AddressID     string
	UserID        string
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	Street        string
	PostalCode    *string
	AddressLabel  string
	Longitude     *decimal.Decimal
	Latitude      *decimal.Decimal
}

// GetAddressesResp 获取地址列表响应
type GetAddressesResp struct {
	Addresses []*AddressDTO
	Total     int
}

// AddressDTO 地址DTO
type AddressDTO struct {
	ID            string
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	Street        string
	PostalCode    *string
	AddressLabel  string
	IsDefault     bool
	Longitude     *decimal.Decimal
	Latitude      *decimal.Decimal
	FullAddress   string
	CreatedAt     int64
	UpdatedAt     int64
}

// AddAddress 添加地址
func (s *Service) AddAddress(ctx context.Context, req *AddAddressReq) (*AddAddressResp, error) {
	// 创建地址实体
	opts := &address.AddressOptions{
		PostalCode:   req.PostalCode,
		AddressLabel: req.AddressLabel,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
	}

	addr, err := address.CreateAddress(
		req.UserID,
		req.ReceiverName,
		req.ReceiverPhone,
		req.Province,
		req.City,
		req.District,
		req.Street,
		opts,
	)
	if err != nil {
		return nil, err
	}

	// 添加地址
	if err := s.addressDomain.AddAddress(ctx, addr); err != nil {
		return nil, err
	}

	return &AddAddressResp{
		Address: s.toDTO(addr),
	}, nil
}

// UpdateAddress 更新地址
func (s *Service) UpdateAddress(ctx context.Context, req *UpdateAddressReq) (*AddAddressResp, error) {
	// 获取现有地址
	addr, err := s.addressRepo.GetByID(ctx, req.AddressID)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errors.New("地址不存在")
	}

	// 检查权限
	if addr.UserID != req.UserID {
		return nil, errors.New("无权限操作此地址")
	}

	// 更新地址信息
	opts := &address.AddressOptions{
		PostalCode:   req.PostalCode,
		AddressLabel: req.AddressLabel,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
	}

	if err := addr.Update(
		req.ReceiverName,
		req.ReceiverPhone,
		req.Province,
		req.City,
		req.District,
		req.Street,
		opts,
	); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.addressDomain.UpdateAddress(ctx, addr); err != nil {
		return nil, err
	}

	return &AddAddressResp{
		Address: s.toDTO(addr),
	}, nil
}

// DeleteAddress 删除地址
func (s *Service) DeleteAddress(ctx context.Context, userID, addressID string) error {
	return s.addressDomain.DeleteAddress(ctx, userID, addressID)
}

// GetAddresses 获取地址列表
func (s *Service) GetAddresses(ctx context.Context, userID string) (*GetAddressesResp, error) {
	addresses, err := s.addressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*AddressDTO, len(addresses))
	for i, addr := range addresses {
		dtos[i] = s.toDTO(addr)
	}

	return &GetAddressesResp{
		Addresses: dtos,
		Total:     len(dtos),
	}, nil
}

// GetAddress 获取地址详情
func (s *Service) GetAddress(ctx context.Context, userID, addressID string) (*AddressDTO, error) {
	addr, err := s.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errors.New("地址不存在")
	}

	// 检查权限
	if addr.UserID != userID {
		return nil, errors.New("无权限访问此地址")
	}

	return s.toDTO(addr), nil
}

// GetDefaultAddress 获取默认地址
func (s *Service) GetDefaultAddress(ctx context.Context, userID string) (*AddressDTO, error) {
	addr, err := s.addressRepo.GetDefaultByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, errors.New("用户没有默认地址")
	}

	return s.toDTO(addr), nil
}

// SetDefaultAddress 设置默认地址
func (s *Service) SetDefaultAddress(ctx context.Context, userID, addressID string) error {
	return s.addressDomain.SetDefaultAddress(ctx, userID, addressID)
}

// toDTO 转换为DTO
func (s *Service) toDTO(addr *address.Address) *AddressDTO {
	dto := &AddressDTO{
		ID:            addr.ID,
		ReceiverName:  addr.ReceiverName,
		ReceiverPhone: addr.ReceiverPhone,
		Province:      addr.Province,
		City:          addr.City,
		District:      addr.District,
		Street:        addr.Street,
		PostalCode:    addr.PostalCode,
		AddressLabel:  addr.AddressLabel,
		IsDefault:     addr.IsDefault,
		Longitude:     addr.Longitude,
		Latitude:      addr.Latitude,
		FullAddress:   addr.GetFullAddress(),
		CreatedAt:     addr.CreatedAt.Unix(),
		UpdatedAt:     addr.UpdatedAt.Unix(),
	}

	return dto
}
