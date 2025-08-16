package user

import (
	"context"
	"errors"
	"fmt"
)

// DomainService 用户领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建用户领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

// RegisterUser 用户注册领域服务
func (s *DomainService) RegisterUser(ctx context.Context, username, email, password string, phone *string) (*User, error) {
	// 检查邮箱是否已存在
	if exists, err := s.repo.ExistsByEmail(ctx, email); err != nil {
		return nil, fmt.Errorf("检查邮箱是否存在失败: %w", err)
	} else if exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 检查手机号是否已存在（如果提供）
	if phone != nil && *phone != "" {
		if exists, err := s.repo.ExistsByPhone(ctx, *phone); err != nil {
			return nil, fmt.Errorf("检查手机号是否存在失败: %w", err)
		} else if exists {
			return nil, errors.New("手机号已被注册")
		}
	}

	// 生成用户名（使用邮箱前缀）
	generatedUsername := s.generateUsernameFromEmail(email)
	if username == "" {
		username = generatedUsername
	}

	// 确保用户名唯一
	uniqueUsername, err := s.ensureUniqueUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("生成唯一用户名失败: %w", err)
	}

	// 创建用户实体
	user, err := CreateUser(uniqueUsername, email, password, phone)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// AuthenticateUser 用户认证领域服务
func (s *DomainService) AuthenticateUser(ctx context.Context, account, password string) (*User, error) {
	// 根据账号查找用户
	user, err := s.repo.FindByAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := user.VerifyPassword(password); err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// ChangeUserEmail 更换用户邮箱领域服务
func (s *DomainService) ChangeUserEmail(ctx context.Context, userID, newEmail string) error {
	// 检查新邮箱是否已存在
	if exists, err := s.repo.ExistsByEmail(ctx, newEmail); err != nil {
		return fmt.Errorf("检查邮箱是否存在失败: %w", err)
	} else if exists {
		return errors.New("邮箱已被使用")
	}

	// 查找用户
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("查找用户失败: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 更换邮箱
	if err := user.ChangeEmail(newEmail); err != nil {
		return fmt.Errorf("更换邮箱失败: %w", err)
	}

	// 保存用户
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("保存用户失败: %w", err)
	}

	return nil
}

// ResetUserPassword 重置用户密码领域服务
func (s *DomainService) ResetUserPassword(ctx context.Context, email, newPassword string) error {
	// 查找用户
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("查找用户失败: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 更新密码
	if err := user.UpdatePassword(newPassword); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 保存用户
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("保存用户失败: %w", err)
	}

	return nil
}

// generateUsernameFromEmail 从邮箱生成用户名
func (s *DomainService) generateUsernameFromEmail(email string) string {
	// 取邮箱@前面的部分作为用户名
	at := len(email)
	for i, r := range email {
		if r == '@' {
			at = i
			break
		}
	}
	return email[:at]
}

// ensureUniqueUsername 确保用户名唯一
func (s *DomainService) ensureUniqueUsername(ctx context.Context, username string) (string, error) {
	originalUsername := username
	counter := 1

	for {
		exists, err := s.repo.ExistsByUsername(ctx, username)
		if err != nil {
			return "", err
		}
		if !exists {
			return username, nil
		}
		username = fmt.Sprintf("%s%d", originalUsername, counter)
		counter++

		// 防止无限循环
		if counter > 1000 {
			return "", errors.New("无法生成唯一用户名")
		}
	}
}
