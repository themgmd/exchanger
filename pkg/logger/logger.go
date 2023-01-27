package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(filePath string) *zap.SugaredLogger {
	writeSyncer := getLogWriter(filePath)
	encoder := getEncoder()

	syncer := zap.CombineWriteSyncers(os.Stdout, writeSyncer)
	core := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())

	return logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    2,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}

	return zapcore.AddSync(lumberJackLogger)
}
