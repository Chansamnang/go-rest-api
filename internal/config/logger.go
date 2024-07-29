package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

func LoggerInit() {
	var err error
	Logger, err = setupLogger()
	if err != nil {
		panic(err)
	}
}

func setupLogger() (*zap.Logger, error) {
	infoLogFile, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	errorLogFile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	// Create WriteSyncers
	infoWriteSyncer := zapcore.AddSync(infoLogFile)
	errorWriteSyncer := zapcore.AddSync(errorLogFile)

	// Create zapcore encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = "" // Disable stack traces in the logs

	// Create cores for different log levels
	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		infoWriteSyncer,
		zap.InfoLevel,
	)
	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		errorWriteSyncer,
		zap.ErrorLevel,
	)

	// Combine cores into a single logger
	logger := zap.New(zapcore.NewTee(infoCore, errorCore), zap.AddCaller())
	return logger, nil
}
