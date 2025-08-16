package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/people257/poor-guy-shop/user-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/user-service/internal/domain/user"
	"gorm.io/gorm"
)

var _ user.Repository = (*UserRepository)(nil)

// UserRepository 用户仓储实现
type UserRepository struct {
	db        *gorm.DB
	q         *query.Query
	converter *user.Converter
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB, q *query.Query, converter *user.Converter) user.Repository {
	return &UserRepository{
		db:        db,
		q:         q,
		converter: converter,
	}
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	userModel, err := r.q.User.WithContext(ctx).Where(r.q.User.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToDomain(userModel), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	userModel, err := r.q.User.WithContext(ctx).Where(r.q.User.Username.Eq(username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToDomain(userModel), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	userModel, err := r.q.User.WithContext(ctx).Where(r.q.User.Email.Eq(email)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToDomain(userModel), nil
}

// FindByPhone 根据手机号查找用户
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*user.User, error) {
	userModel, err := r.q.User.WithContext(ctx).Where(r.q.User.PhoneNumber.Eq(phone)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToDomain(userModel), nil
}

// FindByAccount 根据账号（用户名、邮箱或手机号）查找用户
func (r *UserRepository) FindByAccount(ctx context.Context, account string) (*user.User, error) {
	// 先尝试按用户名查找
	if u, err := r.FindByUsername(ctx, account); err != nil {
		return nil, err
	} else if u != nil {
		return u, nil
	}

	// 判断是否为邮箱格式
	if strings.Contains(account, "@") {
		return r.FindByEmail(ctx, account)
	}

	// 判断是否为手机号格式（简单判断）
	if len(account) == 11 && account[0] == '1' {
		return r.FindByPhone(ctx, account)
	}

	return nil, nil
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	userModel := r.converter.ToModel(u)
	err := r.q.User.WithContext(ctx).Create(userModel)
	if err != nil {
		return err
	}
	u.ID = userModel.ID
	return nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	userModel := r.converter.ToModel(u)
	_, err := r.q.User.WithContext(ctx).Where(r.q.User.ID.Eq(u.ID)).Updates(userModel)
	return err
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.q.User.WithContext(ctx).Where(r.q.User.Username.Eq(username)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.q.User.WithContext(ctx).Where(r.q.User.Email.Eq(email)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByPhone 检查手机号是否存在
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	count, err := r.q.User.WithContext(ctx).Where(r.q.User.PhoneNumber.Eq(phone)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
