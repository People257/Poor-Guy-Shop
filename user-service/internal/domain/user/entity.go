package user

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户领域实体
type User struct {
	ID           string
	Username     string
	PasswordHash *string
	PhoneNumber  *string
	Email        *string
	AvatarURL    *string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserStatus 用户状态
type UserStatus int16

const (
	UserStatusNormal   UserStatus = 1 // 正常
	UserStatusDisabled UserStatus = 2 // 已禁用
)

// CreateUser 创建新用户
func CreateUser(username, email, password string, phone *string) (*User, error) {
	// 验证用户名
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}

	// 验证邮箱
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	// 验证密码
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 验证手机号（如果提供）
	if phone != nil && *phone != "" {
		if err := ValidatePhoneNumber(*phone); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	return &User{
		Username:     username,
		PasswordHash: &hashedPassword,
		Email:        &email,
		PhoneNumber:  phone,
		Status:       UserStatusNormal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) error {
	if u.PasswordHash == nil {
		return errors.New("用户未设置密码")
	}
	return bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password))
}

// IsActive 检查用户是否为活跃状态
func (u *User) IsActive() bool {
	return u.Status == UserStatusNormal
}

// ChangeEmail 更换邮箱
func (u *User) ChangeEmail(newEmail string) error {
	if err := ValidateEmail(newEmail); err != nil {
		return err
	}
	u.Email = &newEmail
	u.UpdatedAt = time.Now()
	return nil
}

// UpdatePassword 更新密码
func (u *User) UpdatePassword(newPassword string) error {
	if err := ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.PasswordHash = &hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return errors.New("用户名长度必须在3-50字符之间")
	}

	// 用户名只能包含字母、数字、下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少为8位")
	}

	if len(password) > 128 {
		return errors.New("密码长度不能超过128位")
	}

	// 至少包含一个字母和一个数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return errors.New("密码必须包含至少一个字母和一个数字")
	}

	return nil
}

// ValidatePhoneNumber 验证手机号格式
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return nil // 手机号可以为空
	}

	// 简单的手机号验证，支持中国大陆手机号
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("手机号格式不正确")
	}

	return nil
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
