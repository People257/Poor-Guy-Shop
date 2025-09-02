package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/user-service/gen/gen/model"
	"github.com/people257/poor-guy-shop/user-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/user-service/internal/domain/address"
)

// addressRepository 地址仓储实现
type addressRepository struct {
	db *gorm.DB
	q  *query.Query
}

// NewAddressRepository 创建地址仓储
func NewAddressRepository(db *gorm.DB) address.Repository {
	return &addressRepository{
		db: db,
		q:  query.Use(db),
	}
}

// Create 创建地址
func (r *addressRepository) Create(ctx context.Context, addr *address.Address) error {
	userAddress := r.toModel(addr)

	if err := r.q.UserAddress.WithContext(ctx).Create(userAddress); err != nil {
		return err
	}

	// 更新实体ID
	addr.ID = userAddress.ID
	return nil
}

// Update 更新地址
func (r *addressRepository) Update(ctx context.Context, addr *address.Address) error {
	userAddress := r.toModel(addr)

	_, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.ID.Eq(addr.ID)).
		Updates(userAddress)

	return err
}

// Delete 删除地址
func (r *addressRepository) Delete(ctx context.Context, id string) error {
	_, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.ID.Eq(id)).
		Delete()

	return err
}

// GetByID 根据ID获取地址
func (r *addressRepository) GetByID(ctx context.Context, id string) (*address.Address, error) {
	userAddress, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.ID.Eq(id)).
		First()

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(userAddress), nil
}

// GetByUserID 获取用户的所有地址
func (r *addressRepository) GetByUserID(ctx context.Context, userID string) ([]*address.Address, error) {
	userAddresses, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.UserID.Eq(userID)).
		Order(r.q.UserAddress.IsDefault.Desc(), r.q.UserAddress.CreatedAt.Desc()).
		Find()

	if err != nil {
		return nil, err
	}

	addresses := make([]*address.Address, len(userAddresses))
	for i, ua := range userAddresses {
		addresses[i] = r.toDomain(ua)
	}

	return addresses, nil
}

// GetDefaultByUserID 获取用户默认地址
func (r *addressRepository) GetDefaultByUserID(ctx context.Context, userID string) (*address.Address, error) {
	userAddress, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.UserID.Eq(userID), r.q.UserAddress.IsDefault.Is(true)).
		First()

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(userAddress), nil
}

// CountByUserID 统计用户地址数量
func (r *addressRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	count, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.UserID.Eq(userID)).
		Count()

	return int(count), err
}

// UnsetAllDefault 取消用户所有默认地址
func (r *addressRepository) UnsetAllDefault(ctx context.Context, userID string) error {
	_, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.UserID.Eq(userID)).
		Update(r.q.UserAddress.IsDefault, false)

	return err
}

// SetDefault 设置默认地址
func (r *addressRepository) SetDefault(ctx context.Context, userID, addressID string) error {
	_, err := r.q.UserAddress.WithContext(ctx).
		Where(r.q.UserAddress.UserID.Eq(userID), r.q.UserAddress.ID.Eq(addressID)).
		Update(r.q.UserAddress.IsDefault, true)

	return err
}

// toModel 转换为数据库模型
func (r *addressRepository) toModel(addr *address.Address) *model.UserAddress {
	return &model.UserAddress{
		ID:            addr.ID,
		UserID:        addr.UserID,
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
		CreatedAt:     addr.CreatedAt,
		UpdatedAt:     addr.UpdatedAt,
	}
}

// toDomain 转换为领域实体
func (r *addressRepository) toDomain(ua *model.UserAddress) *address.Address {
	return &address.Address{
		ID:            ua.ID,
		UserID:        ua.UserID,
		ReceiverName:  ua.ReceiverName,
		ReceiverPhone: ua.ReceiverPhone,
		Province:      ua.Province,
		City:          ua.City,
		District:      ua.District,
		Street:        ua.Street,
		PostalCode:    ua.PostalCode,
		AddressLabel:  ua.AddressLabel,
		IsDefault:     ua.IsDefault,
		Longitude:     ua.Longitude,
		Latitude:      ua.Latitude,
		CreatedAt:     ua.CreatedAt,
		UpdatedAt:     ua.UpdatedAt,
	}
}
