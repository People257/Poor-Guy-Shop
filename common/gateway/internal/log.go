package internal

import (
	"context"
	"fmt"
	"github.com/people257/poor-guy-shop/common/gateway/config"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewZapLogger 创建新的 Zap 日志记录器
func NewZapLogger(cfg *config.ServerConfig,
	logCfg *config.LogConfig,
	observabilityCfg *config.ObservabilityConfig,
	loggerProvider log.LoggerProvider,
) (*zap.Logger, func()) {
	// 根据环境创建编码器配置
	encoderConfig := createEncoderConfig(cfg.Env)

	// 创建日志核心数组,用于支持多输出
	var cores []zapcore.Core

	// 如果启用了文件日志,添加文件输出核心
	if logCfg.File.Enable {
		cores = append(cores, createFileCore(logCfg, encoderConfig))
	}

	// 如果启用了控制台日志,添加控制台输出核心
	if logCfg.Console.Enable {
		cores = append(cores, createConsoleCore(logCfg, encoderConfig))
	}

	// 如果启用可观测性,添加可观测性输出核心
	if observabilityCfg.Log.Enable {
		cores = append(cores, otelzap.NewCore(cfg.Name, otelzap.WithLoggerProvider(loggerProvider)))
	}

	// 使用 Tee 创建多核心日志记录器
	core := zapcore.NewTee(cores...)

	// 创建日志记录器并配置选项
	logger := zap.New(
		core,
		zap.AddCaller(),                       // 添加调用者信息
		zap.AddStacktrace(zapcore.ErrorLevel), // Error 级别及以上添加堆栈跟踪
	)

	// 替换全局日志记录器
	zap.ReplaceGlobals(logger)

	return logger, createCleanupFunc(logger)
}

// createEncoderConfig 创建日志编码器配置
func createEncoderConfig(env string) zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 生产环境使用小写日志级别,其他环境使用彩色大写日志级别
	if env == config.EnvProd {
		cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	} else {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return cfg
}

// createFileCore 创建文件输出核心
// 配置日志文件的路径、轮转策略等
func createFileCore(cfg *config.LogConfig, encConfig zapcore.EncoderConfig) zapcore.Core {
	// 确保日志目录存在
	logDir := cfg.File.Directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create log directory: %v", err))
	}

	// 配置日志轮转
	rotator := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, cfg.File.Name), // 日志文件路径
		MaxSize:    cfg.File.MaxSize,                     // 单个文件最大尺寸(MB)
		MaxAge:     cfg.File.MaxAge,                      // 保留天数
		MaxBackups: cfg.File.MaxBackups,                  // 保留的旧文件数量
		Compress:   cfg.File.Compress,                    // 是否压缩
		LocalTime:  cfg.File.LocalTime,                   // 使用本地时间
	}

	lvl, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		panic(err)
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encConfig), // JSON 格式编码器
		zapcore.AddSync(rotator),          // 添加同步写入器
		lvl,                               // 设置日志级别
	)
}

// createConsoleCore 创建控制台输出核心
// 支持 JSON 或普通文本格式
func createConsoleCore(cfg *config.LogConfig, encConfig zapcore.EncoderConfig) zapcore.Core {
	var encoder zapcore.Encoder
	if cfg.Console.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encConfig)
	}

	lvl, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		panic(err)
	}
	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout), // 标准输出
		lvl,                        // 设置日志级别
	)
}

// createCleanupFunc 创建日志清理函数
// 在服务关闭时同步缓存的日志
func createCleanupFunc(logger *zap.Logger) func() {
	return func() {
		_ = logger.Sync()
	}
}

func NewLoggerProvider(
	exporter sdklog.Exporter,
	serverCfg *config.ServerConfig,
	cfg *config.ObservabilityConfig,

) (log.LoggerProvider, func()) {
	cleanUp := func() {}

	if cfg.Log.Enable {
		res := resource.NewSchemaless(
			semconv.ServiceName(serverCfg.Name),
			semconv.DeploymentEnvironmentName(serverCfg.Env),
		)

		p := sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
			sdklog.WithResource(res),
		)
		cleanUp = func() {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()
			if err := p.Shutdown(ctx); err != nil {
				zap.L().Error("failed to shutdown trace provider", zap.Error(err))
			}
		}
		return p, cleanUp
	} else {
		return noop.NewLoggerProvider(), cleanUp
	}
}

func NewLogExporter(ctx context.Context, cfg *config.ObservabilityConfig) (sdklog.Exporter, func()) {
	if !cfg.Log.Enable {
		return nil, func() {}
	}
	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.Address),
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithCompressor("gzip"),
		otlploggrpc.WithHeaders(cfg.OTLPHeaders),
	)
	if err != nil {
		panic(err)
	}
	cleanUp := func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()
		if err := exporter.Shutdown(ctx); err != nil {
			zap.L().Error("failed to shutdown log exporter", zap.Error(err))
		}
	}
	return exporter, cleanUp
}
