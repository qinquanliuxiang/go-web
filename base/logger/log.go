package logger

import (
	"context"
	"os"
	"qqlx/base/conf"
	"qqlx/base/constant"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Caller() *zap.SugaredLogger {
	return zap.S().WithOptions(zap.AddCaller())
}

func InitLogger() {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.FullCallerEncoder,
	}
	var encoder zapcore.Encoder
	encoder = zapcore.NewJSONEncoder(config)

	writer := zapcore.AddSync(os.Stdout)
	var logLevelStr = conf.GetLogLevel()
	var logLevel zapcore.Level
	switch logLevelStr {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "err":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.InfoLevel
	}
	core := zapcore.NewCore(encoder, writer, logLevel)

	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
	zap.S().Infof("log initialization successful, log level: %s", logLevelStr)
}

func WithContext(ctx context.Context, addCaller bool) *zap.SugaredLogger {
	if addCaller {
		lg := zap.S().WithOptions(zap.AddCaller())
		if traceID := ctx.Value(constant.TraceID).(string); traceID != "" {
			return lg.With(constant.TraceID, traceID)
		}
		return lg
	}
	return zap.S()
}
