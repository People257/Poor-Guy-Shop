package captcha

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/people257/poor-guy-shop/user-service/internal/domain/auth"
)

var _ auth.CaptchaService = (*EmailCaptchaService)(nil)

// EmailCaptchaService 邮箱验证码服务实现
type EmailCaptchaService struct {
	// TODO: 添加邮件发送服务和Redis缓存
}

// NewEmailCaptchaService 创建邮箱验证码服务
func NewEmailCaptchaService() auth.CaptchaService {
	return &EmailCaptchaService{}
}

// SendEmailOTP 发送邮箱验证码
func (s *EmailCaptchaService) SendEmailOTP(ctx context.Context, email string, purpose string) error {
	// 生成6位数字验证码
	otp := s.generateOTP()

	// TODO: 存储验证码到Redis，设置过期时间（如5分钟）
	// key := fmt.Sprintf("email_otp:%s:%s", email, purpose)
	// redis.Set(key, otp, 5*time.Minute)

	// TODO: 发送邮件
	// 这里需要集成邮件发送服务（如阿里云、腾讯云等）

	// 临时日志输出，实际应该发送邮件
	fmt.Printf("发送验证码到邮箱 %s (用途: %s): %s\n", email, purpose, otp)

	return nil
}

// VerifyEmailOTP 验证邮箱验证码
func (s *EmailCaptchaService) VerifyEmailOTP(ctx context.Context, email, otp, purpose string) error {
	// TODO: 从Redis获取存储的验证码进行验证
	// key := fmt.Sprintf("email_otp:%s:%s", email, purpose)
	// storedOTP := redis.Get(key)
	// if storedOTP != otp {
	//     return errors.New("验证码错误或已过期")
	// }
	// redis.Del(key) // 验证成功后删除验证码

	// 临时实现：固定验证码 123456
	if otp != "123456" {
		return fmt.Errorf("验证码错误")
	}

	return nil
}

// generateOTP 生成6位数字验证码
func (s *EmailCaptchaService) generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
