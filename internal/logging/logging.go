package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomLogger struct {
	Log *zap.SugaredLogger
}

type LoggerOptions struct {
	LogFilePath os.File // if you set flag -lf then, enable console + log file
	LogLevel    string  // by the moment only supported info or debug
	Sugared     bool    // Enable sugared logger
}

func NewLogger(options *LoggerOptions) *CustomLogger {
	return &CustomLogger{Log: logger(options)}
}

func logger(options *LoggerOptions) *zap.SugaredLogger {
	//config := zap.NewProductionConfig()
	//
	//config.EncoderConfig = zapcore.EncoderConfig{
	//	EncodeTime:    zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000"),
	//	TimeKey:       "timestamp",
	//	StacktraceKey: "", // to hide stacktrace info
	//	EncodeLevel:   zapcore.CapitalColorLevelEncoder,
	//}
	
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")
	config.TimeKey = "timestamp"
	config.StacktraceKey = "" // to hide stacktrace info
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Pending to add file logger support
	//logFile, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//writer := zapcore.AddSync(logFile)

	defaultLogLevel, _ := zapcore.ParseLevel(options.LogLevel)

	core := zapcore.NewTee(
		//zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger.Sugar()
}
