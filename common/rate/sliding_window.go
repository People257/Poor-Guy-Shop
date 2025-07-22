package rate

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math"
	"time"
)

// <key>:<window>
const limitKeyPattern = "{%s}:%d"

var _ Limiter = (*SlidingWindowLimiter)(nil)

type SlidingWindowLimiter struct {
	rdb redis.UniversalClient
}

func NewSlidingWindowLimiter(rdb redis.UniversalClient) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{rdb: rdb}
}

func (s *SlidingWindowLimiter) Allow(ctx context.Context, key string, limit Limit) bool {
	return s.AllowN(ctx, key, limit, 1)
}

func (s *SlidingWindowLimiter) AllowN(ctx context.Context, key string, limit Limit, n int) bool {
	prevWindowCountKey := s.getPrevWindowCountKey(key, limit)
	currentWindowCountKey := s.getCurrentWindowCountKey(key, limit)

	pipeline := s.rdb.Pipeline()

	prevWindowCountCmd := pipeline.Get(ctx, prevWindowCountKey)
	currentWindowCountCmd := pipeline.Get(ctx, currentWindowCountKey)

	// process error when using result
	_, _ = pipeline.Exec(ctx)

	prevWindowCount, err := prevWindowCountCmd.Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.L().Error("failed to get prev window count", zap.Error(err))
		return true // fail open
	}
	currentWindowCount, err := currentWindowCountCmd.Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		zap.L().Error("failed to get current window count", zap.Error(err))
		return true // fail open
	}

	if currentWindowCount > int64(limit.limit) {
		return false
	}

	// Reference: https://blog.cloudflare.com/counting-things-a-lot-of-different-things/
	current := time.Now()
	currentWindowStart := time.Now().Truncate(limit.windowSize)
	prev := current.Add(-limit.windowSize)

	prevWindowCountWeight := float64(currentWindowStart.Sub(prev)) / float64(limit.windowSize)

	windowCount := float64(prevWindowCount)*prevWindowCountWeight + float64(currentWindowCount)

	if int(math.Round(windowCount))+1 > limit.limit {
		return false
	}

	// 增加当前窗口计数
	pipeline = s.rdb.Pipeline()
	incrCmd := pipeline.IncrBy(ctx, currentWindowCountKey, int64(n))
	expireCmd := pipeline.ExpireNX(ctx, currentWindowCountKey, limit.windowSize*2)

	_, err = pipeline.Exec(ctx)
	err = incrCmd.Err()
	if err != nil {
		zap.L().Error("failed to incr current window count", zap.Error(err))
		return true // fail open
	}
	_, err = expireCmd.Result()
	if err != nil {
		zap.L().Warn("failed to expire current window count", zap.Error(err))
	}

	return true
}

func (s *SlidingWindowLimiter) getCurrentWindowCountKey(key string, limit Limit) string {
	return fmt.Sprintf(limitKeyPattern,
		key,
		time.Now().Truncate(limit.windowSize).UnixMilli())
}

func (s *SlidingWindowLimiter) getPrevWindowCountKey(key string, limit Limit) string {
	return fmt.Sprintf(limitKeyPattern,
		key,
		time.Now().Truncate(limit.windowSize).Add(-limit.windowSize).UnixMilli())
}
