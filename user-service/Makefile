# Makefile for User Service

# Go 相关变量
GO := go
GOFMT := gofmt
GOIMPORTS := goimports
GOLINT := golangci-lint

# 项目信息
PROJECT_NAME := user-service
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_REVISION := $(shell git rev-parse --short HEAD)

# 目录变量
SRC_DIR := .
BUILD_DIR := ./build
PROTO_DIR := ./proto
GEN_DIR := ./gen

# 构建标志
LDFLAGS := -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitRevision=${GIT_REVISION}"

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 开发相关命令
.PHONY: fmt
fmt: ## 格式化代码
	@echo "格式化 Go 代码..."
	@$(GOFMT) -s -w .
	@echo "整理导入包..."
	@$(GOIMPORTS) -w .

.PHONY: lint
lint: ## 运行代码检查
	@echo "运行 golangci-lint..."
	@$(GOLINT) run ./...

.PHONY: vet
vet: ## 运行 go vet
	@echo "运行 go vet..."
	@$(GO) vet ./...

.PHONY: tidy
tidy: ## 整理依赖
	@echo "整理 Go 模块依赖..."
	@$(GO) mod tidy

.PHONY: clean-cache
clean-cache: ## 清理 Go 缓存
	@echo "清理 Go 缓存..."
	@$(GO) clean -modcache
	@$(GO) clean -cache

# 测试相关命令
.PHONY: test
test: ## 运行测试
	@echo "运行单元测试..."
	@$(GO) test -v -race -timeout 30s ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

.PHONY: test-integration
test-integration: ## 运行集成测试
	@echo "运行集成测试..."
	@$(GO) test -v -tags=integration ./tests/...

.PHONY: benchmark
benchmark: ## 运行性能测试
	@echo "运行性能测试..."
	@$(GO) test -bench=. -benchmem ./...

# 代码生成命令
.PHONY: gen
gen: gen-proto gen-gorm gen-wire ## 生成所有代码

.PHONY: gen-proto
gen-proto: ## 生成 protobuf 代码
	@echo "生成 protobuf 代码..."
	@buf generate

.PHONY: gen-gorm
gen-gorm: ## 生成 GORM 代码
	@echo "生成 GORM 查询代码..."
	@$(GO) run ./cmd/gen/main.go

.PHONY: gen-wire
gen-wire: ## 生成 Wire 依赖注入代码
	@echo "生成 Wire 代码..."
	@wire ./cmd/grpc
	@wire ./cmd/gateway

.PHONY: gen-mock
gen-mock: ## 生成 Mock 代码
	@echo "生成 Mock 代码..."
	@$(GO) generate ./...

# 构建相关命令
.PHONY: build
build: build-grpc build-gateway ## 构建所有服务

.PHONY: build-grpc
build-grpc: ## 构建 gRPC 服务
	@echo "构建 gRPC 服务..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/grpc-server ./cmd/grpc

.PHONY: build-gateway
build-gateway: ## 构建 HTTP 网关
	@echo "构建 HTTP 网关..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/gateway-server ./cmd/gateway

.PHONY: build-linux
build-linux: ## 构建 Linux 版本
	@echo "构建 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/grpc-server-linux ./cmd/grpc
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/gateway-server-linux ./cmd/gateway

# 运行相关命令
.PHONY: run-grpc
run-grpc: ## 运行 gRPC 服务
	@echo "启动 gRPC 服务..."
	@$(GO) run ./cmd/grpc

.PHONY: run-gateway
run-gateway: ## 运行 HTTP 网关
	@echo "启动 HTTP 网关..."
	@$(GO) run ./cmd/gateway

.PHONY: run-dev
run-dev: ## 开发模式运行（热重载）
	@echo "开发模式启动..."
	@air

# Docker 相关命令
.PHONY: docker-build
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	@docker build -t $(PROJECT_NAME):$(VERSION) .
	@docker build -t $(PROJECT_NAME):latest .

.PHONY: docker-push
docker-push: ## 推送 Docker 镜像
	@echo "推送 Docker 镜像..."
	@docker push $(PROJECT_NAME):$(VERSION)
	@docker push $(PROJECT_NAME):latest

.PHONY: docker-run
docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	@docker run -p 8080:8080 -p 9090:9090 $(PROJECT_NAME):latest

# 清理相关命令
.PHONY: clean
clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

.PHONY: clean-gen
clean-gen: ## 清理生成的代码
	@echo "清理生成的代码..."
	@rm -rf $(GEN_DIR)/proto
	@rm -rf $(GEN_DIR)/gen
	@find . -name "wire_gen.go" -delete
	@find . -name "*_mock.go" -delete

.PHONY: clean-all
clean-all: clean clean-gen clean-cache ## 清理所有文件

# 数据库相关命令
.PHONY: db-migrate
db-migrate: ## 运行数据库迁移
	@echo "运行数据库迁移..."
	@$(GO) run ./cmd/migrate

.PHONY: db-seed
db-seed: ## 初始化测试数据
	@echo "初始化测试数据..."
	@$(GO) run ./cmd/seed

.PHONY: db-reset
db-reset: ## 重置数据库
	@echo "重置数据库..."
	@$(GO) run ./cmd/migrate -reset

# 质量检查命令
.PHONY: check
check: fmt vet lint test ## 运行所有质量检查

.PHONY: pre-commit
pre-commit: tidy fmt vet lint test ## 提交前检查

.PHONY: ci
ci: tidy fmt vet lint test-coverage ## CI 流水线检查

# 文档相关命令
.PHONY: docs
docs: ## 生成文档
	@echo "生成 API 文档..."
	@$(GO) run ./cmd/docs

.PHONY: docs-serve
docs-serve: ## 启动文档服务器
	@echo "启动文档服务器..."
	@$(GO) run ./cmd/docs -serve

# 工具安装命令
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "安装开发工具..."
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@$(GO) install github.com/google/wire/cmd/wire@latest
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install github.com/cosmtrek/air@latest
	@$(GO) install go.uber.org/mock/mockgen@latest
	@$(GO) install gorm.io/gen/tools/gentool@latest

# 版本相关命令
.PHONY: version
version: ## 显示版本信息
	@echo "项目版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git 版本: $(GIT_REVISION)"

.PHONY: tag
tag: ## 创建版本标签
	@read -p "请输入版本号 (例如: v1.0.0): " version; \
	git tag -a $$version -m "Release $$version"; \
	echo "已创建标签: $$version"

# 性能分析命令
.PHONY: pprof-cpu
pprof-cpu: ## CPU 性能分析
	@echo "启动 CPU 性能分析..."
	@$(GO) tool pprof http://localhost:6060/debug/pprof/profile

.PHONY: pprof-heap
pprof-heap: ## 内存性能分析
	@echo "启动内存性能分析..."
	@$(GO) tool pprof http://localhost:6060/debug/pprof/heap

# 安全检查命令
.PHONY: security
security: ## 运行安全检查
	@echo "运行安全检查..."
	@$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

# 项目初始化命令
.PHONY: init
init: install-tools tidy gen ## 初始化项目

.PHONY: setup
setup: init db-migrate ## 完整项目设置

# 监控和健康检查
.PHONY: health
health: ## 检查服务健康状态
	@echo "检查服务健康状态..."
	@curl -f http://localhost:8080/health || echo "服务未启动"

.PHONY: metrics
metrics: ## 查看服务指标
	@echo "查看服务指标..."
	@curl -s http://localhost:8080/metrics

# 部署相关命令
.PHONY: deploy-dev
deploy-dev: ## 部署到开发环境
	@echo "部署到开发环境..."
	@./scripts/deploy.sh dev

.PHONY: deploy-staging
deploy-staging: ## 部署到测试环境
	@echo "部署到测试环境..."
	@./scripts/deploy.sh staging

.PHONY: deploy-prod
deploy-prod: ## 部署到生产环境
	@echo "部署到生产环境..."
	@./scripts/deploy.sh prod