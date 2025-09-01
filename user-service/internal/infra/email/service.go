package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/people257/poor-guy-shop/user-service/internal/config"
)

// Service 邮件服务接口
type Service interface {
	SendTemplate(ctx context.Context, templateName, to string, data interface{}) error
	SendText(ctx context.Context, to, subject, body string) error
}

// SMTPService SMTP邮件服务实现
type SMTPService struct {
	config *config.EmailConfig
}

// NewSMTPService 创建SMTP邮件服务
func NewSMTPService(cfg *config.EmailConfig) Service {
	return &SMTPService{
		config: cfg,
	}
}

// SendTemplate 使用模板发送邮件
func (s *SMTPService) SendTemplate(ctx context.Context, templateName, to string, data interface{}) error {
	// 获取邮件模板
	tmpl, exists := s.config.Templates[templateName]
	if !exists {
		return fmt.Errorf("邮件模板不存在: %s", templateName)
	}

	// 解析主题模板
	subjectTmpl, err := template.New("subject").Parse(tmpl.Subject)
	if err != nil {
		return fmt.Errorf("解析主题模板失败: %w", err)
	}

	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, data); err != nil {
		return fmt.Errorf("执行主题模板失败: %w", err)
	}

	// 解析内容模板
	bodyTmpl, err := template.New("body").Parse(tmpl.Body)
	if err != nil {
		return fmt.Errorf("解析内容模板失败: %w", err)
	}

	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, data); err != nil {
		return fmt.Errorf("执行内容模板失败: %w", err)
	}

	return s.SendText(ctx, to, subjectBuf.String(), bodyBuf.String())
}

// SendText 发送纯文本邮件
func (s *SMTPService) SendText(ctx context.Context, to, subject, body string) error {
	// 构建符合RFC标准的邮件内容
	message := s.buildRFC5322Message(to, subject, body)

	// SMTP服务器地址
	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)

	// 创建认证
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)

	// 如果使用TLS，使用STARTTLS方式
	if s.config.SMTP.UseTLS {
		err := s.sendWithSTARTTLS(addr, auth, s.config.SMTP.From, []string{to}, []byte(message))
		if err != nil {
			return fmt.Errorf("STARTTLS邮件发送失败 [%s -> %s]: %w", s.config.SMTP.From, to, err)
		}
		return nil
	}

	// 普通SMTP发送
	err := smtp.SendMail(addr, auth, s.config.SMTP.From, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("SMTP邮件发送失败 [%s:%d, %s -> %s]: %w", s.config.SMTP.Host, s.config.SMTP.Port, s.config.SMTP.From, to, err)
	}
	return nil
}

// sendWithSTARTTLS 使用STARTTLS发送邮件（适用于587端口）
func (s *SMTPService) sendWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 创建普通TCP连接
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %w", err)
	}
	defer conn.Quit()

	// 启动TLS
	tlsConfig := &tls.Config{
		ServerName: s.config.SMTP.Host,
	}

	if err = conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("启动TLS失败: %w", err)
	}

	// 认证
	if auth != nil {
		if err = conn.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %w", err)
		}
	}

	// 设置发件人
	if err = conn.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err = conn.Rcpt(recipient); err != nil {
			return fmt.Errorf("设置收件人失败: %w", err)
		}
	}

	// 发送邮件内容
	writer, err := conn.Data()
	if err != nil {
		return fmt.Errorf("获取邮件写入器失败: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("关闭邮件写入器失败: %w", err)
	}

	return nil
}

// buildRFC5322Message 构建符合RFC5322标准的邮件消息
func (s *SMTPService) buildRFC5322Message(to, subject, body string) string {
	var message strings.Builder

	// 必需的头部字段
	message.WriteString(fmt.Sprintf("From: %s\r\n", s.config.SMTP.From))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	// MIME相关头部
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: 8bit\r\n")

	// 可选但推荐的头部
	message.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	message.WriteString("X-Mailer: User-Service-SMTP/1.0\r\n")

	// 消息ID（可选但有助于追踪）
	message.WriteString(fmt.Sprintf("Message-ID: <%d.%s>\r\n", time.Now().UnixNano(), s.config.SMTP.From))

	// 空行分隔头部和正文
	message.WriteString("\r\n")

	// 邮件正文
	message.WriteString(body)

	return message.String()
}
