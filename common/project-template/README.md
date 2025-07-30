# Project Template
后端服务模板

## 从模版创建
1. 安装 `gonew` 工具
    ```shell
    go install golang.org/x/tools/cmd/gonew@latest
    ```
2. 执行 
    ```shell
    gonew cnb.cool/cymirror/ces-services/common/project-template cnb.cool/cymirror/ces-services/xxx-service
    ```

## 工具链
### protobuf toolchain
- [Buf](https://buf.build/docs/tutorials/getting-started-with-buf-cli/)

## 使用
### swagger + api接口 + proto/grpc 定义
```shell
buf generate
```

### 生成数据库模型/查询
#### Gen
```shell
go run ./cmd/gen
```

## 其他
### 目录结构
```
├── README.md     // README
├── api           // api，grpc 接口实现
├── buf.gen.yaml  // buf 生成配置
├── buf.lock      // buf.lock
├── buf.yaml      // buf 配置
├── cmd           // 程序入口（gen 模型 / query 生成，Gateway，gRPC Server）
├── gen           // 生成的代码
├── go.mod        // go.mod
├── go.sum        // go.sum
├── internal      // 内部模块
│        ├── application     // 应用层，负责逻辑编排
│        ├── domain          // 领域层，专注业务逻辑
│        └── infrastructure  // 基础设施层，负责与外部交互
└── proto         // proto定义
```

### VSCode/JetBrains 插件
- Buf

## 参考
- [Buf Cli](https://buf.build/docs/tutorials/getting-started-with-buf-cli/)
- [grpc-ecosystem proto demo](https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto)
- [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [protovalidate example](https://github.com/bufbuild/protovalidate/tree/main/examples)
- [gorm/gen](https://github.com/go-gorm/gen)