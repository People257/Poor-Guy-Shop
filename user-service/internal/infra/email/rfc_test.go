package email

import (
	"strings"
	"testing"

	"github.com/people257/poor-guy-shop/user-service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRFC5322MessageFormat(t *testing.T) {
	// 准备测试配置
	cfg := &config.EmailConfig{
		SMTP: config.SMTPConfig{
			From: "test@example.com",
		},
	}

	service := NewSMTPService(cfg).(*SMTPService)

	// 测试邮件消息构建
	message := service.buildRFC5322Message("recipient@example.com", "Test Subject", "Test Body")

	t.Run("必需头部字段检查", func(t *testing.T) {
		// 检查From头部
		assert.Contains(t, message, "From: test@example.com\r\n", "必须包含From头部")

		// 检查To头部
		assert.Contains(t, message, "To: recipient@example.com\r\n", "必须包含To头部")

		// 检查Subject头部
		assert.Contains(t, message, "Subject: Test Subject\r\n", "必须包含Subject头部")
	})

	t.Run("MIME头部检查", func(t *testing.T) {
		// 检查MIME版本
		assert.Contains(t, message, "MIME-Version: 1.0\r\n", "必须包含MIME版本")

		// 检查Content-Type
		assert.Contains(t, message, "Content-Type: text/plain; charset=UTF-8\r\n", "必须包含Content-Type")

		// 检查Content-Transfer-Encoding
		assert.Contains(t, message, "Content-Transfer-Encoding: 8bit\r\n", "必须包含Content-Transfer-Encoding")
	})

	t.Run("可选头部检查", func(t *testing.T) {
		// 检查Date头部
		assert.Contains(t, message, "Date: ", "应该包含Date头部")

		// 检查X-Mailer头部
		assert.Contains(t, message, "X-Mailer: User-Service-SMTP/1.0\r\n", "应该包含X-Mailer头部")

		// 检查Message-ID头部
		assert.Contains(t, message, "Message-ID: <", "应该包含Message-ID头部")
		assert.Contains(t, message, ".test@example.com>\r\n", "Message-ID应该包含发送方域名")
	})

	t.Run("消息结构检查", func(t *testing.T) {
		// 检查头部和正文分隔符
		assert.Contains(t, message, "\r\n\r\n", "头部和正文之间必须有空行分隔")

		// 检查正文内容
		assert.True(t, strings.HasSuffix(message, "Test Body"), "消息应该以正文内容结尾")
	})

	t.Run("RFC5322合规性检查", func(t *testing.T) {
		lines := strings.Split(message, "\r\n")

		// 检查是否有空行分隔头部和正文
		emptyLineFound := false
		headerEnded := false

		for i, line := range lines {
			if line == "" && !headerEnded {
				emptyLineFound = true
				headerEnded = true
				// 确保正文在空行之后
				if i+1 < len(lines) {
					assert.Equal(t, "Test Body", lines[i+1], "正文应该在空行之后")
				}
				break
			}
		}

		assert.True(t, emptyLineFound, "必须有空行分隔头部和正文")
	})
}

func TestEmailHeaderEncoding(t *testing.T) {
	cfg := &config.EmailConfig{
		SMTP: config.SMTPConfig{
			From: "发送者@example.com", // 测试中文域名
		},
	}

	service := NewSMTPService(cfg).(*SMTPService)

	// 测试包含中文的主题和正文
	message := service.buildRFC5322Message("收件人@example.com", "测试主题", "测试正文内容")

	t.Run("中文内容处理", func(t *testing.T) {
		// 检查UTF-8编码声明
		assert.Contains(t, message, "charset=UTF-8", "必须声明UTF-8编码")

		// 检查中文内容是否正确包含
		assert.Contains(t, message, "测试正文内容", "正文中文内容应该正确包含")
	})
}

func TestQQMailCompatibility(t *testing.T) {
	// 测试QQ邮箱特定要求
	cfg := &config.EmailConfig{
		SMTP: config.SMTPConfig{
			From: "sender@qq.com",
		},
	}

	service := NewSMTPService(cfg).(*SMTPService)
	message := service.buildRFC5322Message("recipient@qq.com", "QQ Mail Test", "QQ Mail Body")

	t.Run("QQ邮箱格式要求", func(t *testing.T) {
		// QQ邮箱要求的关键头部
		requiredHeaders := []string{
			"From: sender@qq.com",
			"To: recipient@qq.com",
			"Subject: QQ Mail Test",
			"MIME-Version: 1.0",
			"Content-Type: text/plain; charset=UTF-8",
		}

		for _, header := range requiredHeaders {
			assert.Contains(t, message, header, "QQ邮箱要求的头部: %s", header)
		}

		// 检查正确的CRLF换行符
		assert.NotContains(t, message, "\n\n", "不应该有单独的LF换行符")
		assert.Contains(t, message, "\r\n\r\n", "应该使用CRLF换行符")
	})
}
