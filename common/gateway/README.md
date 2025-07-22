# Gateway

基于 gRPC Gateway，实现 gRPC => REST

## 注意

- HTTP Header 相关
    - [映射规则](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/runtime#DefaultHeaderMatcher)：Mapping HTTP
      headers with `Grpc-Metadata-` prefix to gRPC metadata (prefixed with `grpcgateway-`)
    - 常用默认映射且***不添加*** `grpcgateway-`  前缀的请求头
        - `X-Forwarded-For`
        - `X-Forwarded-Host`
        - `Authorization`
        - `Content-Type`
    - 常用默认映射且***添加*** `grpcgateway-`  前缀的请求头
        - 详见[源码](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/context.go#L334)
    - 映射完存放在 gRPC Metadata 的 Header key 全部为小写，虽然获取的时候不区分大小写但是会多一次遍历的成本 (遍历判断
      key == ToLower(key))
- Response Body 相关
    - 默认情况下，零值不返回，相当于 json tag 中的 `omitempty`

## 使用

1. 通过 `Buf CLI` 工具，执行 `buf generate` 命令即可生成包括 *.pb.gw.go 在内的代码文件
2. 通过 `gateway.InitializeGateway` 初始化 Gateway （配置文件详见 `config.yaml.example`）
3. 在 RegisterHandler 函数中注册 Handler，添加 `xxxpb.RegisterxxxHandler(context.Background(), gwmux, conn)`
4. 启动 Gateway

```go
ctx := context.TODO()
errGroup, ctx := errgroup.WithContext(ctx)

var cfg *config.GatewayConfig
// 从其他地方读取配置，配置文件见 .config.yaml.example
// cfg = config.Load()

gw, cleanUp := gateway.InitializeGateway(ctx, cfg)

// 注册 xxx 服务，注册 Protobuf 定义的 HTTP 接口
_ = gw.RegisterHandler(func (gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
    var err error
    err = foopb.RegisterFooServiceHandler(context.Background(), gwmux, conn)
    if err != nil {
        panic(err)
    }
    return nil
})

errGroup.Go(func () error {
    return gw.Run(ctx)
})
errGroup.Go(func () error {
    <-ctx.Done()
    cleanUp()
    return nil
})

if err := errGroup.Wait(); err != nil {
    // 处理错误
}
```