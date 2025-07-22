package rate

import "context"

type Limiter interface {
	Allow(ctx context.Context, key string, limit Limit) bool
	AllowN(ctx context.Context, key string, limit Limit, n int) bool
}
