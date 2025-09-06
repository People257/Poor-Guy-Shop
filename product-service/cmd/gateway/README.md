# Gateway

## 使用
1. 使用 `Buf Cli` 生成 Gateway 等代码
2. 前往 `application.go` 中注册到 Mux
3. 添加 `xxxpb.RegisterxxxHandler(context.Background(), gwmux, conn)` 到代码中
4. 启动 Gateway

## 其他
### 添加配置
添加配置一般在 `cmd/gateway/internal/config/config.go` 中
```go
type Config struct {
    config.GatewayConfig `mapstructure:",squash"`
    // 添加配置
    Foo string `mapstructure:"foo"`
}
```
然后在 provider.go 中添加 GetFoo 方法提供给 wire