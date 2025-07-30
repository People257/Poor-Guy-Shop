# gRPC Server

## 使用
1. 使用 `Buf Cli` 生成 ProtoBuf，gRPC 等代码
2. 实现 `xxxpb.xxxServer` 接口
3. 前往 `application.go` 中注册 Server
4. 启动 Server

## 其他
### 添加配置
添加配置一般在 `cmd/grpc/internal/config/config.go` 中
```go
type Config struct {
    config.GrpcServerConfig `mapstructure:",squash"`
    // 添加配置
    Foo string `mapstructure:"foo"`
}
```
然后在 provider.go 中添加 GetFoo 方法提供给 wire