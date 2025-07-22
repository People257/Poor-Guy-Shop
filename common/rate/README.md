# Rate Limiter

基于 Cloudflare Blog 的 Rate Limiter 实现，使用 Redis 作为存储引擎。

## 使用

```go
var limiter rate.Limiter
limiter = rate.NewSlidingWindowLimiter(rdb, 10)

// 窗口为 1 秒，每秒允许 1 次，消耗 1 个令牌
allow := limiter.AllowN(ctx, "key", rate.PerSecond(1), 1)
fmt.Println(allow)

// Allow(ctx, "key") 与 AllowN(ctx, "key", 1) 相同
allow := limiter.Allow(ctx, "key")
fmt.Println(allow)
```

## 参考
- [Cloudflare Blog](https://blog.cloudflare.com/counting-things-a-lot-of-different-things/)