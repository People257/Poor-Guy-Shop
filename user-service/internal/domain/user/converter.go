package user

import (
	"github.com/people257/poor-guy-shop/user-service/gen/gen/model"
	"time"
)

// Converter 用户实体转换器
type Converter struct{}

// NewConverter 创建转换器
func NewConverter() *Converter {
	return &Converter{}
}

// ToModel 将领域实体转换为数据模型
func (c *Converter) ToModel(u *User) *model.User {
	if u == nil {
		return nil
	}

	return &model.User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PhoneNumber:  u.PhoneNumber,
		Email:        u.Email,
		AvatarURL:    u.AvatarURL,
		Status:       int16(u.Status),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// ToDomain 将数据模型转换为领域实体
func (c *Converter) ToDomain(m *model.User) *User {
	if m == nil {
		return nil
	}

	return &User{
		ID:           m.ID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		PhoneNumber:  m.PhoneNumber,
		Email:        m.Email,
		AvatarURL:    m.AvatarURL,
		Status:       UserStatus(m.Status),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToModels 批量将领域实体转换为数据模型
func (c *Converter) ToModels(users []*User) []*model.User {
	if len(users) == 0 {
		return nil
	}

	models := make([]*model.User, len(users))
	for i, u := range users {
		models[i] = c.ToModel(u)
	}
	return models
}

// ToDomains 批量将数据模型转换为领域实体
func (c *Converter) ToDomains(models []*model.User) []*User {
	if len(models) == 0 {
		return nil
	}

	users := make([]*User, len(models))
	for i, m := range models {
		users[i] = c.ToDomain(m)
	}
	return users
}

// NewUserFromModel 从数据模型创建新的用户实体（用于从数据库加载）
func (c *Converter) NewUserFromModel(m *model.User) (*User, error) {
	if m == nil {
		return nil, nil
	}

	// 验证必要字段
	if err := ValidateUsername(m.Username); err != nil {
		return nil, err
	}

	if m.Email != nil && *m.Email != "" {
		if err := ValidateEmail(*m.Email); err != nil {
			return nil, err
		}
	}

	if m.PhoneNumber != nil && *m.PhoneNumber != "" {
		if err := ValidatePhoneNumber(*m.PhoneNumber); err != nil {
			return nil, err
		}
	}

	return &User{
		ID:           m.ID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		PhoneNumber:  m.PhoneNumber,
		Email:        m.Email,
		AvatarURL:    m.AvatarURL,
		Status:       UserStatus(m.Status),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}

// PrepareForCreate 为创建操作准备实体（设置默认值）
func (c *Converter) PrepareForCreate(u *User) *User {
	if u == nil {
		return nil
	}

	now := time.Now()
	prepared := *u // 复制结构体

	// 设置创建和更新时间
	prepared.CreatedAt = now
	prepared.UpdatedAt = now

	// 确保状态正确
	if prepared.Status == 0 {
		prepared.Status = UserStatusNormal
	}

	return &prepared
}

// PrepareForUpdate 为更新操作准备实体
func (c *Converter) PrepareForUpdate(u *User) *User {
	if u == nil {
		return nil
	}

	prepared := *u // 复制结构体
	prepared.UpdatedAt = time.Now()

	return &prepared
}
