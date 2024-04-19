package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var logger *zap.SugaredLogger

func init() {
	logger = NewLogger().Sugar()
}

func NewLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	config := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig:     encoderConfig,
	}

	// 根据环境变量设置日志级别
	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		var level zap.AtomicLevel
		if err := level.UnmarshalText([]byte(lvl)); err == nil {
			config.Level = level
		}
	}

	zLogger, err := config.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	return zLogger
}

// Now, use SugaredLogger methods directly
func Info(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}

func Sync() {
	logger.Sync()
}
