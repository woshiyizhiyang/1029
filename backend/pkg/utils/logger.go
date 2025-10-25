package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger 初始化日志系统
func InitLogger(logPath string) error {
	var err error
	once.Do(func() {
		// 配置日志编码器
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 配置日志级别
		atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

		// 控制台输出
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			atomicLevel,
		)

		// 文件输出
		var fileCore zapcore.Core
		if logPath != "" {
			file, openErr := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if openErr != nil {
				err = openErr
				return
			}
			fileCore = zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(file),
				atomicLevel,
			)
		}

		// 组合核心
		var core zapcore.Core
		if fileCore != nil {
			core = zapcore.NewTee(consoleCore, fileCore)
		} else {
			core = consoleCore
		}

		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})

	return err
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if logger == nil {
		_ = InitLogger("")
	}
	return logger
}

// Info 记录Info级别日志
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Error 记录Error级别日志
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Warn 记录Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Debug 记录Debug级别日志
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Fatal 记录Fatal级别日志
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}
