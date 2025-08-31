package captcha

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/people257/poor-guy-shop/user-service/internal/config"

	"github.com/people257/poor-guy-shop/user-service/internal/domain/auth"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/email"
	"github.com/redis/go-redis/v9"
)

var _ auth.CaptchaService = (*EmailCaptchaService)(nil)

// EmailCaptchaService 邮箱验证码服务实现
type EmailCaptchaService struct {
	emailService email.Service
	redisClient  *redis.Client
	config       *config.CaptchaConfig
}

// NewEmailCaptchaService 创建邮箱验证码服务
func NewEmailCaptchaService(
	emailService email.Service,
	redisClient *redis.Client,
	config *config.CaptchaConfig,
) auth.CaptchaService {
	return &EmailCaptchaService{
		emailService: emailService,
		redisClient:  redisClient,
		config:       config,
	}
}

// SendEmailOTP 发送邮箱验证码
func (s *EmailCaptchaService) SendEmailOTP(ctx context.Context, email string, purpose string) error {
	// 检查配置是否启用
	if !s.config.Email.Enabled {
		return fmt.Errorf("邮箱验证码功能未启用")
	}

	// 检查发送频率限制
	intervalKey := fmt.Sprintf("email_otp_interval:%s:%s", email, purpose)
	exists, err := s.redisClient.Exists(ctx, intervalKey).Result()
	if err != nil {
		return fmt.Errorf("检查发送间隔失败: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("发送太频繁，请等待%d秒后重试", s.config.Email.SendInterval)
	}

	// 检查每日发送限制
	dailyKey := fmt.Sprintf("email_otp_daily:%s:%s", email, time.Now().Format("2006-01-02"))
	dailyCount, err := s.redisClient.Get(ctx, dailyKey).Int()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("检查每日限制失败: %w", err)
	}
	if dailyCount >= s.config.Email.DailyLimit {
		return fmt.Errorf("今日发送次数已达上限(%d次)", s.config.Email.DailyLimit)
	}

	// 生成验证码
	otp := s.generateOTP()

	// 存储验证码到Redis
	otpKey := fmt.Sprintf("email_otp:%s:%s", email, purpose)
	err = s.redisClient.Set(ctx, otpKey, otp, time.Duration(s.config.Email.ExpiresIn)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("存储验证码失败: %w", err)
	}

	// 设置发送间隔限制
	err = s.redisClient.Set(ctx, intervalKey, "1", time.Duration(s.config.Email.SendInterval)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("设置发送间隔失败: %w", err)
	}

	// 增加每日发送计数
	err = s.redisClient.Incr(ctx, dailyKey).Err()
	if err != nil {
		return fmt.Errorf("增加每日计数失败: %w", err)
	}
	// 设置每日计数过期时间为第二天0点
	tomorrow := time.Now().AddDate(0, 0, 1)
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	s.redisClient.ExpireAt(ctx, dailyKey, midnight)

	// 根据用途选择邮件模板
	templateName := "verification_code"
	if purpose == "password_reset" {
		templateName = "password_reset"
	}

	// 发送邮件
	err = s.emailService.SendTemplate(ctx, templateName, email, map[string]interface{}{
		"Code": otp,
	})
	if err != nil {
		// 如果邮件发送失败，清理Redis中的验证码
		s.redisClient.Del(ctx, otpKey)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// VerifyEmailOTP 验证邮箱验证码
func (s *EmailCaptchaService) VerifyEmailOTP(ctx context.Context, email, otp, purpose string) error {
	// 检查配置是否启用
	if !s.config.Email.Enabled {
		return fmt.Errorf("邮箱验证码功能未启用")
	}

	// 从Redis获取存储的验证码
	otpKey := fmt.Sprintf("email_otp:%s:%s", email, purpose)
	storedOTP, err := s.redisClient.Get(ctx, otpKey).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("验证码不存在或已过期")
		}
		return fmt.Errorf("获取验证码失败: %w", err)
	}

	// 验证验证码
	if storedOTP != otp {
		return fmt.Errorf("验证码错误")
	}

	// 验证成功后删除验证码（一次性使用）
	err = s.redisClient.Del(ctx, otpKey).Err()
	if err != nil {
		// 删除失败只记录日志，不影响验证结果
		fmt.Printf("删除验证码失败: %v\n", err)
	}

	return nil
}

// generateOTP 生成数字验证码
func (s *EmailCaptchaService) generateOTP() string {
	codeLength := s.config.Email.CodeLength
	if codeLength <= 0 {
		codeLength = 6 // 默认6位
	}

	// 计算最大值
	max := 1
	for i := 0; i < codeLength; i++ {
		max *= 10
	}

	rand.Seed(time.Now().UnixNano())
	format := fmt.Sprintf("%%0%dd", codeLength)
	return fmt.Sprintf(format, rand.Intn(max))
}
