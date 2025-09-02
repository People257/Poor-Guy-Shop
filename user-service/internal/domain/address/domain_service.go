package address

import (
	"context"
	"errors"
)

const MaxAddressCount = 20

// DomainService 地址领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建地址领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

// AddAddress 添加地址
func (s *DomainService) AddAddress(ctx context.Context, address *Address) error {
	// 检查地址数量限制
	count, err := s.repo.CountByUserID(ctx, address.UserID)
	if err != nil {
		return err
	}
	if count >= MaxAddressCount {
		return errors.New("地址数量已达上限，最多支持20个地址")
	}

	// 如果是第一个地址，自动设为默认
	if count == 0 {
		address.SetAsDefault()
	}

	return s.repo.Create(ctx, address)
}

// UpdateAddress 更新地址
func (s *DomainService) UpdateAddress(ctx context.Context, address *Address) error {
	// 检查地址是否存在
	existing, err := s.repo.GetByID(ctx, address.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("地址不存在")
	}

	// 检查权限
	if existing.UserID != address.UserID {
		return errors.New("无权限操作此地址")
	}

	return s.repo.Update(ctx, address)
}

// DeleteAddress 删除地址
func (s *DomainService) DeleteAddress(ctx context.Context, userID, addressID string) error {
	// 检查地址是否存在
	address, err := s.repo.GetByID(ctx, addressID)
	if err != nil {
		return err
	}
	if address == nil {
		return errors.New("地址不存在")
	}

	// 检查权限
	if address.UserID != userID {
		return errors.New("无权限操作此地址")
	}

	// 如果删除的是默认地址，需要设置新的默认地址
	if address.IsDefault {
		addresses, err := s.repo.GetByUserID(ctx, userID)
		if err != nil {
			return err
		}

		// 删除地址
		if err := s.repo.Delete(ctx, addressID); err != nil {
			return err
		}

		// 如果还有其他地址，将最新的设为默认
		for _, addr := range addresses {
			if addr.ID != addressID {
				addr.SetAsDefault()
				return s.repo.Update(ctx, addr)
			}
		}
	} else {
		return s.repo.Delete(ctx, addressID)
	}

	return nil
}

// SetDefaultAddress 设置默认地址
func (s *DomainService) SetDefaultAddress(ctx context.Context, userID, addressID string) error {
	// 检查地址是否存在
	address, err := s.repo.GetByID(ctx, addressID)
	if err != nil {
		return err
	}
	if address == nil {
		return errors.New("地址不存在")
	}

	// 检查权限
	if address.UserID != userID {
		return errors.New("无权限操作此地址")
	}

	// 如果已经是默认地址，直接返回
	if address.IsDefault {
		return nil
	}

	// 取消所有默认地址
	if err := s.repo.UnsetAllDefault(ctx, userID); err != nil {
		return err
	}

	// 设置新的默认地址
	return s.repo.SetDefault(ctx, userID, addressID)
}
