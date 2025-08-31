package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"text/template"

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
	// 构建邮件内容
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		to, subject, body)

	// SMTP服务器地址
	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)

	// 创建认证
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)

	// 如果使用TLS，需要特殊处理
	if s.config.SMTP.UseTLS {
		return s.sendWithTLS(addr, auth, s.config.SMTP.From, []string{to}, []byte(message))
	}

	// 普通SMTP发送
	return smtp.SendMail(addr, auth, s.config.SMTP.From, []string{to}, []byte(message))
}

// sendWithTLS 使用TLS发送邮件
func (s *SMTPService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 创建TLS连接
	tlsConfig := &tls.Config{
		ServerName: s.config.SMTP.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS连接失败: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.config.SMTP.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Quit()

	// 认证
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %w", err)
		}
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("设置收件人失败: %w", err)
		}
	}

	// 发送邮件内容
	writer, err := client.Data()
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
