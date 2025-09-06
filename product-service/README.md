# Product Service

商品管理微服务，基于DDD架构设计，提供分类、品牌、商品的完整管理功能。

## 功能特性

- **分类管理**：支持多级分类、分类树管理
- **品牌管理**：品牌信息的CRUD操作
- **商品管理**：商品和SKU的完整管理
- **搜索功能**：关键词搜索、条件筛选
- **双协议支持**：gRPC + REST API（gRPC-Gateway）

## 架构设计

```
product-service/
├── api/                    # API层（gRPC处理器）
├── cmd/                    # 应用入口
│   ├── grpc/              # gRPC服务
│   ├── gateway/           # HTTP网关
│   └── gen/               # 代码生成工具
├── internal/              # 内部代码
│   ├── domain/            # 领域层
│   ├── application/       # 应用服务层
│   └── infra/            # 基础设施层
├── proto/                 # Proto定义
├── gen/                   # 生成的代码
└── migrations/            # 数据库迁移
```

## 数据库注意事项

### PostgreSQL vs MySQL 语法差异

本项目使用PostgreSQL数据库。在创建表时，注意以下语法差异：

**❌ MySQL风格（不支持）：**
```sql
CREATE TABLE categories (
    name VARCHAR(100) NOT NULL COMMENT '分类名称'
);
```

**✅ PostgreSQL风格（正确）：**
```sql
CREATE TABLE categories (
    name VARCHAR(100) NOT NULL
);

-- 单独添加注释
COMMENT ON TABLE categories IS '商品分类表';
COMMENT ON COLUMN categories.name IS '分类名称';
```

### 数据库迁移

1. 确保PostgreSQL服务已启动
2. 创建数据库：
   ```sql
   CREATE DATABASE product_service;
   ```

3. 运行迁移：
   ```bash
   # 使用migrate工具
   migrate -path migrations -database "postgresql://username:password@localhost/product_service?sslmode=disable" up
   
   # 或直接执行SQL文件
   psql -U username -d product_service -f migrations/001_create_product_tables.sql
   ```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 生成代码

```bash
# 生成Proto代码
make gen-proto

# 生成GORM模型（需要先创建数据库）
make gen-gorm

# 生成Wire依赖注入代码
make gen-wire
```

### 3. 配置文件

复制配置文件模板：
```bash
cp cmd/grpc/etc/config.yaml.example cmd/grpc/etc/config.yaml
```

修改数据库连接信息：
```yaml
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your_password"
  dbname: "product_service"
```

### 4. 运行服务

```bash
# 运行gRPC服务
make run-grpc

# 运行HTTP网关（另一个终端）
make run-gateway
```

## API文档

服务启动后，可以通过以下方式访问API：

- **gRPC**: `localhost:9001`
- **HTTP**: `localhost:8080`
- **Swagger文档**: `http://localhost:8080/swagger/`

### 主要接口

#### 分类管理
- `POST /api/v1/categories` - 创建分类
- `PUT /api/v1/categories/{id}` - 更新分类
- `DELETE /api/v1/categories/{id}` - 删除分类
- `GET /api/v1/categories/{id}` - 获取分类详情
- `GET /api/v1/categories` - 获取分类列表
- `GET /api/v1/categories/tree` - 获取分类树

#### 品牌管理
- `POST /api/v1/brands` - 创建品牌
- `PUT /api/v1/brands/{id}` - 更新品牌
- `DELETE /api/v1/brands/{id}` - 删除品牌
- `GET /api/v1/brands/{id}` - 获取品牌详情
- `GET /api/v1/brands` - 获取品牌列表

#### 商品管理
- `POST /api/v1/products` - 创建商品
- `PUT /api/v1/products/{id}` - 更新商品
- `DELETE /api/v1/products/{id}` - 删除商品
- `GET /api/v1/products/{id}` - 获取商品详情
- `GET /api/v1/products` - 获取商品列表
- `GET /api/v1/products/search` - 搜索商品

## 开发指南

### 代码生成

```bash
# 完整代码生成流程
make gen
```

### 测试

```bash
# 运行测试
make test

# 生成测试覆盖率报告
make test-coverage
```

### 代码检查

```bash
# 代码格式化
make fmt

# 代码检查
make lint
```

## 部署

### Docker部署

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### 环境变量

- `CONFIG_FILE`: 配置文件路径（默认: `etc/config.yaml`）
- `LOG_LEVEL`: 日志级别（默认: `info`）

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查PostgreSQL服务是否启动
   - 验证数据库连接信息是否正确

2. **Proto编译失败**
   - 确保安装了buf工具：`go install github.com/bufbuild/buf/cmd/buf@latest`

3. **Wire生成失败**
   - 确保安装了wire工具：`go install github.com/google/wire/cmd/wire@latest`

4. **GORM生成失败**
   - 确保数据库表已创建
   - 检查数据库连接配置

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License
