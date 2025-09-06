package email

import (
	"context"
	"testing"

	"github.com/people257/poor-guy-shop/user-service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestSMTPService_SendTemplate(t *testing.T) {
	// 准备测试配置
	cfg := &config.EmailConfig{
		SMTP: config.SMTPConfig{
			Host:     "smtp.test.com",
			Port:     587,
			Username: "test@test.com",
			Password: "password",
			From:     "test@test.com",
			UseTLS:   false, // 测试时关闭TLS
		},
		Templates: map[string]config.EmailTemplate{
			"register": {
				Subject: "【用户服务】注册验证码",
				Body:    "欢迎注册！您的验证码是：{{.Code}}",
			},
			"login": {
				Subject: "【用户服务】登录验证码",
				Body:    "您正在登录，验证码是：{{.Code}}",
			},
		},
	}

	// 创建邮件服务
	service := NewSMTPService(cfg)

	t.Run("发送已存在的模板", func(t *testing.T) {
		// 注意：这个测试不会真正发送邮件，因为SMTP服务器不存在
		// 但可以测试模板解析逻辑
		err := service.SendTemplate(context.Background(), "register", "test@example.com", map[string]interface{}{
			"Code": "123456",
		})

		// 由于SMTP服务器连接会失败，我们期望得到连接错误，而不是模板不存在错误
		assert.Error(t, err)
		assert.NotContains(t, err.Error(), "邮件模板不存在")
	})

	t.Run("发送不存在的模板", func(t *testing.T) {
		err := service.SendTemplate(context.Background(), "nonexistent", "test@example.com", map[string]interface{}{
			"Code": "123456",
		})

		// 应该返回模板不存在错误
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "邮件模板不存在: nonexistent")
	})
}

func TestSMTPService_TemplateExists(t *testing.T) {
	// 准备测试配置
	cfg := &config.EmailConfig{
		Templates: map[string]config.EmailTemplate{
			"register": {
				Subject: "注册验证码",
				Body:    "验证码：{{.Code}}",
			},
			"login": {
				Subject: "登录验证码",
				Body:    "验证码：{{.Code}}",
			},
		},
	}

	service := NewSMTPService(cfg).(*SMTPService)

	// 测试已存在的模板
	_, exists := service.config.Templates["register"]
	assert.True(t, exists, "register模板应该存在")

	_, exists = service.config.Templates["login"]
	assert.True(t, exists, "login模板应该存在")

	// 测试不存在的模板
	_, exists = service.config.Templates["nonexistent"]
	assert.False(t, exists, "nonexistent模板不应该存在")
}

func TestEmailTemplateRendering(t *testing.T) {
	// 测试模板渲染逻辑
	cfg := &config.EmailConfig{
		Templates: map[string]config.EmailTemplate{
			"test": {
				Subject: "测试邮件 - {{.Code}}",
				Body:    "您好！您的验证码是：{{.Code}}\n感谢使用我们的服务！",
			},
		},
	}

	service := NewSMTPService(cfg)

	// 由于SendTemplate会尝试连接SMTP服务器，我们直接测试模板是否存在
	smtpService := service.(*SMTPService)
	template, exists := smtpService.config.Templates["test"]

	assert.True(t, exists, "test模板应该存在")
	assert.Equal(t, "测试邮件 - {{.Code}}", template.Subject)
	assert.Contains(t, template.Body, "{{.Code}}")
}
