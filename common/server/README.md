# gRPC Server

## 使用

1. 使用 `Buf Cli` 生成 ProtoBuf，gRPC 等代码
2. 实现 `xxxpb.xxxServer` 接口
3. 在 `RegisterServer` 中注册到 gRPC Server
4. 启动 Server

```go
ctx := context.TODO()
errGroup, ctx := errgroup.WithContext(ctx)

var cfg *config.GrpcServerConfig
// 从其他地方读取配置，配置文件见 .config.yaml.example
// cfg = config.Load()

s, cleanUp := server.InitializeServer(ctx, cfg)

// 注册 xxx Service Server
s.RegisterServer(func(s *grpc.Server) {
    foopb.RegisterFooServiceServer(s, fooServer)
})

errGroup.Go(func() error {
	return s.Run(ctx)
})
errGroup.Go(func() error {
	<-ctx.Done()
	cleanUp()
	return nil
})

if err := errGroup.Wait(); err != nil {
	// 处理错误
}
```