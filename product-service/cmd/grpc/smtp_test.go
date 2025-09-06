// SMTP连接测试工具
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/internal/config"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/email"
)

func main() {
	fmt.Println("🔍 SMTP连接测试工具")

	// 加载配置
	cfg := config.MustLoad("etc/config.yaml")

	// 转换为内部配置
	emailConfig := ProvideInternalEmailConfig(cfg)

	fmt.Printf("📧 SMTP配置信息:\n")
	fmt.Printf("  服务器: %s:%d\n", emailConfig.SMTP.Host, emailConfig.SMTP.Port)
	fmt.Printf("  用户名: %s\n", emailConfig.SMTP.Username)
	fmt.Printf("  发件人: %s\n", emailConfig.SMTP.From)
	fmt.Printf("  使用TLS: %t\n", emailConfig.SMTP.UseTLS)
	fmt.Printf("  模板数量: %d\n\n", len(emailConfig.Templates))

	// 创建邮件服务
	emailService := email.NewSMTPService(emailConfig)

	// 测试发送邮件（发送给自己）
	testEmail := emailConfig.SMTP.From
	fmt.Printf("🧪 测试发送邮件到: %s\n", testEmail)

	ctx := context.Background()
	err := emailService.SendTemplate(ctx, "register", testEmail, map[string]interface{}{
		"Code": "123456",
	})

	if err != nil {
		log.Printf("❌ 邮件发送失败: %v\n", err)

		// 详细的错误诊断
		fmt.Println("\n🔧 错误诊断:")
		if contains(err.Error(), "short response") {
			fmt.Println("  • 'short response' 通常表示:")
			fmt.Println("    - SMTP服务器突然关闭连接")
			fmt.Println("    - 认证信息错误")
			fmt.Println("    - 服务器不支持当前的TLS配置")
			fmt.Println("    - 网络连接问题")
		}

		if contains(err.Error(), "TLS") || contains(err.Error(), "STARTTLS") {
			fmt.Println("  • TLS相关问题:")
			fmt.Println("    - 检查端口587是否正确（通常需要STARTTLS）")
			fmt.Println("    - 端口465通常需要SSL/TLS")
			fmt.Println("    - 端口25通常是明文连接")
		}

		if contains(err.Error(), "认证失败") || contains(err.Error(), "Auth") {
			fmt.Println("  • 认证问题:")
			fmt.Println("    - 检查用户名和密码是否正确")
			fmt.Println("    - QQ邮箱需要使用授权码，不是登录密码")
			fmt.Println("    - 确保邮箱已开启SMTP服务")
		}

		fmt.Println("\n💡 建议解决方案:")
		fmt.Println("  1. 确认QQ邮箱SMTP服务已开启")
		fmt.Println("  2. 使用授权码替代登录密码")
		fmt.Println("  3. 确认网络连接正常")
		fmt.Println("  4. 检查防火墙设置")

	} else {
		fmt.Println("✅ 邮件发送成功！")
		fmt.Println("📨 请检查收件箱确认邮件是否收到")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
