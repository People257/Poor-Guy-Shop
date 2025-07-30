.PHONY: all grpc gateway gen model wire clean help

# 默认配置文件路径
CONFIG_PATH ?= cmd/grpc/etc/config.yaml
GATEWAY_CONFIG_PATH ?= cmd/gateway/etc/config.yaml
GEN_CONFIG_PATH ?= cmd/gen/etc/config.yaml

# 默认目标
all: gen wire

# 运行 gRPC 服务
grpc:
	@echo "Starting gRPC server..."
	go run ./cmd/grpc -f $(CONFIG_PATH)

# 运行 Gateway 服务
gateway:
	@echo "Starting Gateway server..."
	go run ./cmd/gateway -f $(GATEWAY_CONFIG_PATH)

# 生成 Proto 相关代码
gen:
	@echo "Generating Proto files..."
	buf generate

# 生成数据库模型
model:
	@echo "Generating database models..."
	go run ./cmd/gen -f $(GEN_CONFIG_PATH)

# 生成 wire 依赖注入代码
wire:
	@echo "Generating wire dependency injection code..."
	cd cmd/grpc && wire
	cd cmd/gateway && wire
	cd common/gateway && wire
	cd common/server && wire

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  grpc    - Run the gRPC server"
	@echo "  gateway - Run the Gateway server"
	@echo "  gen     - Generate Proto related code"
	@echo "  model   - Generate database models"
	@echo "  wire    - Generate wire dependency injection code"
	@echo "  clean   - Clean generated files"
	@echo "  all     - Generate all code (default)"
	@echo ""
	@echo "Configuration:"
	@echo "  CONFIG_PATH         - gRPC server config path (default: etc/config.yaml)"
	@echo "  GATEWAY_CONFIG_PATH - Gateway server config path (default: etc/config.yaml)" 
	@echo "  GEN_CONFIG_PATH     - Database model generator config path (default: cmd/gen/etc/config.yaml)" 