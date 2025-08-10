package auth

import "time"

// Config 配置
type Config struct {
	Secret           string        // jwt 密钥
	ExpireDuration   time.Duration // jwt 过期时间
	RefreshThreshold time.Duration // 刷新阈值，当剩余时间小于此值时刷新令牌
}
