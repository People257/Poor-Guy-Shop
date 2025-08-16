package auth

import "errors"

// 认证领域错误定义
var (
	// 令牌相关错误
	ErrInvalidToken     = errors.New("无效的令牌")
	ErrTokenExpired     = errors.New("令牌已过期")
	ErrInvalidClaims    = errors.New("无效的令牌声明")
	ErrInvalidUserID    = errors.New("无效的用户ID")
	ErrInvalidExpiresAt = errors.New("无效的过期时间")

	// 验证码相关错误
	ErrInvalidCaptcha  = errors.New("验证码错误")
	ErrCaptchaExpired  = errors.New("验证码已过期")
	ErrCaptchaNotFound = errors.New("验证码不存在")
	ErrTooManyAttempts = errors.New("验证码尝试次数过多")
	ErrCaptchaNotSent  = errors.New("验证码发送失败")

	// 认证相关错误
	ErrUserNotFound       = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("用户已被禁用")
	ErrUserLocked         = errors.New("用户已被锁定")

	// 授权相关错误
	ErrAccessDenied      = errors.New("访问被拒绝")
	ErrInsufficientScope = errors.New("权限不足")
	ErrResourceNotFound  = errors.New("资源不存在")
)
