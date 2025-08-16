package user

import "errors"

// 用户领域错误定义
var (
	// 用户相关错误
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUserAlreadyExists = errors.New("用户已存在")
	ErrUserDisabled      = errors.New("用户已被禁用")
	ErrUserLocked        = errors.New("用户已被锁定")

	// 用户名相关错误
	ErrInvalidUsername       = errors.New("用户名格式不正确")
	ErrUsernameTooShort      = errors.New("用户名长度不能少于3位")
	ErrUsernameTooLong       = errors.New("用户名长度不能超过50位")
	ErrUsernameAlreadyExists = errors.New("用户名已存在")
	ErrUsernameRequired      = errors.New("用户名不能为空")

	// 邮箱相关错误
	ErrInvalidEmail       = errors.New("邮箱格式不正确")
	ErrEmailRequired      = errors.New("邮箱不能为空")
	ErrEmailAlreadyExists = errors.New("邮箱已被注册")

	// 密码相关错误
	ErrInvalidPassword  = errors.New("密码格式不正确")
	ErrPasswordTooShort = errors.New("密码长度不能少于8位")
	ErrPasswordTooLong  = errors.New("密码长度不能超过128位")
	ErrPasswordTooWeak  = errors.New("密码强度不够")
	ErrPasswordRequired = errors.New("密码不能为空")
	ErrPasswordNotSet   = errors.New("用户未设置密码")
	ErrWrongPassword    = errors.New("密码错误")

	// 手机号相关错误
	ErrInvalidPhoneNumber       = errors.New("手机号格式不正确")
	ErrPhoneNumberAlreadyExists = errors.New("手机号已被注册")

	// 状态相关错误
	ErrInvalidUserStatus = errors.New("无效的用户状态")

	// 操作相关错误
	ErrCannotCreateUser   = errors.New("无法创建用户")
	ErrCannotUpdateUser   = errors.New("无法更新用户")
	ErrCannotDeleteUser   = errors.New("无法删除用户")
	ErrOperationForbidden = errors.New("操作被禁止")
)
