package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type CustomLogger struct {
	Log *zap.SugaredLogger
}

type LoggerOptions struct {
	//LogFilePath string // if you set flag -lf then, enable console + log file
	LogLevel string // by the moment only supported info or debug
	Sugar    bool   // Enable sugared logger
}

func NewLogger(options *LoggerOptions) *CustomLogger {
	return &CustomLogger{Log: logger(options)}
}

//func logger(options *LoggerOptions) *zap.SugaredLogger {
//	config := zap.NewProductionEncoderConfig()
//
//	//config.EncoderConfig = zapcore.EncoderConfig{
//	//	EncodeTime:    zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000"),
//	//	TimeKey:       "timestamp",
//	//	StacktraceKey: "", // to hide stacktrace info
//	//	EncodeLevel:   zapcore.CapitalColorLevelEncoder,
//	//}
//
//	config.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")
//	config.TimeKey = "timestamp"
//	config.StacktraceKey = "" // to hide stacktrace info
//	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
//
//	defaultLogLevel, _ := zapcore.ParseLevel(options.LogLevel)
//	config.SetLevel(defaultLogLevel)
//
//	consoleEncoder := zapcore.NewConsoleEncoder(config)
//
//	core := zapcore.NewTee(
//		//zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
//		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
//	)
//
//	if options.LogFilePath != "" {
//		config.OutputPaths = []string{options.LogFilePath, "stderr"}
//	}
//
//	//customLogger, err := config.Build()
//
//	customLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
//
//	defer customLogger.Sync()
//
//	return customLogger.Sugar()
//}

func logger(options *LoggerOptions) *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()

	//config.EncoderConfig = zapcore.EncoderConfig{
	//	EncodeTime:    zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000"),
	//	TimeKey:       "timestamp",
	//	StacktraceKey: "", // to hide stacktrace info
	//	EncodeLevel:   zapcore.CapitalColorLevelEncoder,
	//}

	config.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")
	config.TimeKey = "timestamp"
	config.StacktraceKey = ""
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	defaultLogLevel, _ := zapcore.ParseLevel(options.LogLevel)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Pending to add file logger support
	//fileEncoder := zapcore.NewJSONEncoder(config)
	//logFile, _ := os.OpenFile(options.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		//zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger.Sugar()
}
