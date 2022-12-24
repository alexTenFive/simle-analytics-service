package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Init(d bool) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := zap.InfoLevel
	if d {
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	l := zap.New(core)

	return l.Sugar()
}
