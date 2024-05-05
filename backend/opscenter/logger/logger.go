package logger

import (
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggers sync.Map // 使用 sync.Map 存储和管理不同用户的日志记录器实例
)

// getLogger 创建或返回一个给定用户 ID 的日志记录器
func getLogger(userID string) *zap.Logger {
	if userID == "" {
		userID = "default"
	}
	logPath := filepath.Join("logs", userID+".log")

	if logger, ok := loggers.Load(logPath); ok {
		return logger.(*zap.Logger)
	}

	// 创建一个新的日志记录器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "消息"
	encoderConfig.TimeKey = "时间"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      false,
		DisableCaller:    true,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{logPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	newLogger, err := config.Build()
	if err != nil {
		panic("cannot create logger: " + err.Error())
	}

	loggers.Store(logPath, newLogger)
	return newLogger
}

func Info(userID, message string, fields ...interface{}) {
	getLogger(userID).Sugar().Infow(message, fields...)
}

func Error(userID, message string, fields ...interface{}) {
	getLogger(userID).Sugar().Errorw(message, fields...)
}
