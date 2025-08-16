# OSS 对象存储服务 (简化版)

Poor Guy Shop 的对象存储服务，专注于核心功能：文件上传、下载、安全校验和基本文件管理。

## 功能特性

### 🚀 核心功能
- **文件上传**: 单文件上传，支持多种文件类型
- **文件下载**: 安全的临时URL下载
- **文件管理**: 文件列表查询、删除操作
- **安全校验**: JWT认证，基于所有者的访问控制

### 🛡️ 安全特性
- JWT Token认证
- 基于文件所有者的权限控制
- Public/Private可见性设置
- 基本访问日志记录

### 📁 支持的文件类型
- **图片**: JPEG, PNG, GIF, WebP
- **文档**: PDF, DOC, DOCX, TXT
- **文件分类**: avatar(头像), product(商品), document(文档)

## 快速开始

### 环境要求
- Go 1.24+
- PostgreSQL 12+
- Redis 6+

### 安装依赖
```bash
go mod tidy
```

### 数据库初始化
```bash
# 执行数据库迁移 (简化版)
psql -U postgres -d poor_guy_shop -f migrations/001_create_simple_oss_tables.sql
```

### 配置文件
复制配置示例文件并修改：
```bash
cp cmd/gateway/etc/config.yaml.example cmd/gateway/etc/config.yaml
cp cmd/grpc/etc/config.yaml.example cmd/grpc/etc/config.yaml
```

### 生成代码
```bash
# 生成protobuf代码和swagger文档
buf generate

# 生成数据库模型和查询代码
go run ./cmd/gen
```

### 启动服务
```bash
# 启动gRPC服务
go run ./cmd/grpc

# 启动HTTP网关(另一个终端)
go run ./cmd/gateway
```

服务启动后：
- HTTP API: http://localhost:8080
- gRPC服务: localhost:8081
- Swagger文档: http://localhost:8080/swagger/

## 项目结构

```
oss-infra/
├── README.md                        // 项目说明
├── OSS_SIMPLE_DESIGN.md            // 简化版系统设计文档
├── OSS_SIMPLE_API_GUIDE.md         // 简化版API接口文档
├── migrations/                      // 数据库迁移文件
│   └── 001_create_simple_oss_tables.sql
├── proto/                          // Protocol Buffers定义
│   └── oss/
│       ├── file/                   // 文件服务
│       └── common/                 // 公共定义
├── cmd/                            // 程序入口
│   ├── gateway/                    // HTTP网关
│   ├── grpc/                      // gRPC服务
│   └── gen/                       // 代码生成工具
├── internal/                       // 内部实现
│   ├── application/               // 应用层
│   ├── domain/                   // 领域层
│   └── infra/                    // 基础设施层
├── gen/                           // 生成的代码
│   ├── gen/                      // GORM生成的模型
│   ├── proto/                    // protobuf生成的代码
│   └── swagger/                  // swagger文档
└── api/                          // API接口实现
```

## API 接口

### 核心接口 (简化版)
- `POST /v1/oss/file/upload` - 上传文件
- `GET /v1/oss/file/{id}/download-url` - 获取下载URL
- `GET /v1/oss/file/{id}` - 获取文件信息
- `GET /v1/oss/file/list` - 获取文件列表
- `DELETE /v1/oss/file/{id}` - 删除文件

详细的API文档请查看 [OSS_SIMPLE_API_GUIDE.md](./OSS_SIMPLE_API_GUIDE.md)

## 数据库表结构 (简化版)

### 核心表
- `files` - 文件信息表
- `file_access_logs` - 简单访问日志表

详细的表结构设计请查看 [OSS_SIMPLE_DESIGN.md](./OSS_SIMPLE_DESIGN.md)

## 使用示例

### JavaScript客户端 (简化版)
```javascript
const ossService = new SimpleOSSService('http://localhost:8080', 'your-token');

// 上传文件
const result = await ossService.uploadFile(file, 'avatar', 'private');

// 获取下载URL
const downloadUrl = await ossService.getDownloadUrl(fileId);

// 获取文件列表
const files = await ossService.listFiles('avatar', 1, 10);
```

### Go客户端
```go
client := filepb.NewFileServiceClient(conn)

// 上传文件
resp, err := client.UploadFile(ctx, &filepb.UploadFileReq{
    FileData: fileData,
    Filename: "example.jpg",
    Category: "avatar",
    Visibility: "private",
})
```

## 配置说明

### 基本配置 (简化版)
```yaml
server:
  name: "oss-service"
  port: 8080

storage:
  root_path: "/data/oss"
  url_prefix: "http://localhost:8080/files"
  max_file_size: 52428800  # 50MB

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "poor_guy_shop"

jwt:
  secret: "your-jwt-secret-key"
```

## 开发指南 (简化版)

### 扩展文件类型支持
1. 更新 `AllowedMimeTypes` 配置
2. 添加文件类型验证逻辑
3. 更新前端文件选择器

### 自定义权限规则
1. 修改 `CheckAccess` 函数
2. 添加新的权限检查逻辑
3. 更新相关测试用例

## 监控和运维 (简化版)

### 基本监控
- 存储空间使用率
- API响应时间
- 上传/下载成功率

### 日志管理
- 访问日志记录
- 错误日志监控

### 备份策略
- 数据库定期备份
- 文件存储备份

## 部署

### Docker部署
```bash
# 构建镜像
docker build -t oss-service .

# 运行服务
docker run -d -p 8080:8080 -p 8081:8081 oss-service
```

### Kubernetes部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: oss-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: oss-service
  template:
    spec:
      containers:
      - name: oss-service
        image: oss-service:latest
        ports:
        - containerPort: 8080
        - containerPort: 8081
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目地址: https://github.com/people257/poor-guy-shop
- 问题反馈: https://github.com/people257/poor-guy-shop/issues
- 邮箱: dev@poorguyshop.com

---

*OSS对象存储服务(简化版) - 为Poor Guy Shop提供简洁、安全、易维护的文件存储解决方案*