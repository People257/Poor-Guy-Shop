package internal

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/natefinch/lumberjack"
	"github.com/people257/poor-guy-shop/common/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

func NewZapLogger(
	serverCfg *config.ServerConfig,
	logCfg *config.LogConfig,
) (*zap.Logger, func()) {
	// 创建 zap 日志器
	encoderConfig := createEncoderConfig(serverCfg.Env)

	var cores []zapcore.Core

	if logCfg.Console.Enable {
		cores = append(cores, createConsoleCore(logCfg, encoderConfig))
	}

	if logCfg.File.Enable {
		cores = append(cores, createFileCore(logCfg, encoderConfig))
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	zap.ReplaceGlobals(logger)

	return logger, createCleanupFunc(logger)
}

func createEncoderConfig(env string) zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 默认使用小写级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 格式的时间
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 只显示文件名和行号
	}

	if env != config.EnvProd {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	return cfg
}

func createFileCore(cfg *config.LogConfig, encConfig zapcore.EncoderConfig) zapcore.Core {
	logDir := cfg.File.Directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create log directory: %v", err))
	}

	rotator := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, cfg.File.Name), // 日志文件路径
		MaxSize:    int(cfg.File.MaxSize),                // 单个文件最大尺寸 (MB)
		MaxAge:     int(cfg.File.MaxAge),                 // 文件最长保留天数
		MaxBackups: int(cfg.File.MaxBackups),             // 最多保留的旧文件数量
		Compress:   cfg.File.Compress,                    // 是否压缩旧文件
		LocalTime:  cfg.File.LocalTime,                   // 是否使用本地时间
	}

	// 解析配置中的日志级别字符串，例如 "info", "debug", "error"。
	logLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		// 如果配置的级别无效，默认使用 info 级别。
		logLevel = zapcore.InfoLevel
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encConfig), // 文件日志总是使用 JSON 格式，方便机器解析。
		zapcore.AddSync(rotator),          // 将 lumberjack rotator 作为写入目标。
		logLevel,                          // 设置这个 Core 处理的最低日志级别。
	)
}

func createConsoleCore(cfg *config.LogConfig, encConfig zapcore.EncoderConfig) zapcore.Core {
	var encoder zapcore.Encoder
	if cfg.Console.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encConfig)
	}

	logLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		// 如果配置的级别无效，默认使用 info 级别。
		logLevel = zapcore.InfoLevel
	}

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		logLevel,
	)

}

// createCleanupFunc 创建日志清理函数，在日志关闭时执行，保证日志文件写入完成
func createCleanupFunc(logger *zap.Logger) func() {
	return func() {
		_ = logger.Sync()
	}
}

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}
		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("invalid log level: %v", lvl))
		}
	})
}
