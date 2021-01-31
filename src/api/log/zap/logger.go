package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log var
	Log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"stdout"},
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	var err error
	Log, err = logConfig.Build()
	if err != nil {
		panic(err)
	}
}

// Field func
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Debug func
func Debug(msg string, tags ...zap.Field) {
	Log.Debug(msg, tags...)
	Log.Sync()
}

// Info func
func Info(msg string, tags ...zap.Field) {
	Log.Info(msg, tags...)
	Log.Sync()
}

// Error func
func Error(msg string, err error, tags ...zap.Field) {
	msg = fmt.Sprintf("%s - ERROR - %v", msg, err)
	Log.Error(msg, tags...)
	Log.Sync()
}
