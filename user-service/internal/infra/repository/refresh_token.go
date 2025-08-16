package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/people257/poor-guy-shop/user-service/internal/domain/auth"
	"github.com/redis/go-redis/v9"
)

var _ auth.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

// RefreshTokenRepository 刷新令牌仓储实现
type RefreshTokenRepository struct {
	redis *redis.Client
}

// NewRefreshTokenRepository 创建刷新令牌仓储
func NewRefreshTokenRepository(redis *redis.Client) auth.RefreshTokenRepository {
	return &RefreshTokenRepository{
		redis: redis,
	}
}

// Store 存储刷新令牌
func (r *RefreshTokenRepository) Store(ctx context.Context, token *auth.RefreshToken) error {
	// 序列化令牌数据
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("序列化刷新令牌失败: %w", err)
	}

	// 存储到 Redis，设置过期时间
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return errors.New("刷新令牌已过期")
	}

	// 使用令牌值作为 key
	key := r.getTokenKey(token.Token)
	if err := r.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("存储刷新令牌到Redis失败: %w", err)
	}

	// 同时建立用户ID到令牌的映射，用于快速删除用户的所有令牌
	userKey := r.getUserTokensKey(token.UserID)
	if err := r.redis.SAdd(ctx, userKey, token.Token).Err(); err != nil {
		return fmt.Errorf("添加用户令牌映射失败: %w", err)
	}

	// 设置用户令牌集合的过期时间
	if err := r.redis.Expire(ctx, userKey, ttl).Err(); err != nil {
		return fmt.Errorf("设置用户令牌集合过期时间失败: %w", err)
	}

	return nil
}

// FindByToken 根据令牌值查找
func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*auth.RefreshToken, error) {
	key := r.getTokenKey(token)
	data, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // 令牌不存在
		}
		return nil, fmt.Errorf("从Redis获取刷新令牌失败: %w", err)
	}

	var refreshToken auth.RefreshToken
	if err := json.Unmarshal([]byte(data), &refreshToken); err != nil {
		return nil, fmt.Errorf("反序列化刷新令牌失败: %w", err)
	}

	return &refreshToken, nil
}

// MarkAsUsed 标记令牌为已使用
func (r *RefreshTokenRepository) MarkAsUsed(ctx context.Context, tokenID string) error {
	// 由于我们使用令牌值作为key，需要先找到令牌
	// 这里简化处理，直接删除令牌（一次性使用）
	// 在实际实现中，你可能需要更新令牌状态
	return nil // 暂时不实现，通过删除令牌来实现一次性使用
}

// DeleteByUserID 删除用户的所有刷新令牌
func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	userKey := r.getUserTokensKey(userID)

	// 获取用户的所有令牌
	tokens, err := r.redis.SMembers(ctx, userKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil // 用户没有令牌
		}
		return fmt.Errorf("获取用户令牌列表失败: %w", err)
	}

	// 删除所有令牌
	if len(tokens) > 0 {
		keys := make([]string, len(tokens))
		for i, token := range tokens {
			keys[i] = r.getTokenKey(token)
		}

		if err := r.redis.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("删除用户令牌失败: %w", err)
		}
	}

	// 删除用户令牌集合
	if err := r.redis.Del(ctx, userKey).Err(); err != nil {
		return fmt.Errorf("删除用户令牌集合失败: %w", err)
	}

	return nil
}

// DeleteExpired 删除过期的令牌
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	// Redis 会自动删除过期的 key，这里不需要特别处理
	// 如果需要清理用户令牌集合中的过期令牌引用，可以在这里实现
	return nil
}

// DeleteToken 删除特定令牌
func (r *RefreshTokenRepository) DeleteToken(ctx context.Context, token string, userID string) error {
	// 删除令牌
	key := r.getTokenKey(token)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("删除刷新令牌失败: %w", err)
	}

	// 从用户令牌集合中移除
	userKey := r.getUserTokensKey(userID)
	if err := r.redis.SRem(ctx, userKey, token).Err(); err != nil {
		return fmt.Errorf("从用户令牌集合中移除令牌失败: %w", err)
	}

	return nil
}

// getTokenKey 获取令牌的Redis key
func (r *RefreshTokenRepository) getTokenKey(token string) string {
	return fmt.Sprintf("refresh_token:%s", token)
}

// getUserTokensKey 获取用户令牌集合的Redis key
func (r *RefreshTokenRepository) getUserTokensKey(userID string) string {
	return fmt.Sprintf("user_tokens:%s", userID)
}
