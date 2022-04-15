package log

import (
	"cs/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var logger *Logger

// Logger 日志类
type Logger struct {
	Zap *zap.Logger
	ZapSugar *zap.SugaredLogger
	AppLevel *zap.AtomicLevel
}

// GetLogger 获取日志实例
func GetLogger() *Logger {
	return logger
}

// SetAppLogLevel 设置日志级别
func (l *Logger) SetAppLogLevel(level string) {
	if l.AppLevel == nil {
		return
	}
	newLevel := getLogLevel(level)
	if l.AppLevel.String() != newLevel.String() {
		l.AppLevel.SetLevel(newLevel)
		l.Zap.Info("set new log level: " + newLevel.String())
	}
}

func Init() {
	logConf := config.MustUse("log")
	logPath := logConf.MustString("path", "./log")
	logLevel := logConf.MustString("level", "info")
	appLogName := logConf.MustString("appLogName", "app")
	errLogName := logConf.MustString("errLogName", "err")
	maxSize := logConf.MustInt("maxSize", 500)  // 日志轮换前文件最大大小
	maxAge := logConf.MustInt("maxAge", 30)  // 旧日志最长保留时间
	maxBackups := logConf.MustInt("maxBackups", 100)  // 保留旧日志文件最大数量
	compress := logConf.MustBool("compress", false)  // 是否使用gzip压缩日志文件
	showConsole := logConf.MustBool("showConsole", false)

	appFile := filepath.Join(filepath.FromSlash(logPath), appLogName)
	errFile := filepath.Join(filepath.FromSlash(logPath), errLogName)

	appLogLevel := zap.NewAtomicLevelAt(getLogLevel(logLevel))
	errLogLevel := zap.NewAtomicLevelAt(getLogLevel("error"))

	appLumberjackHook := lumberjackHook(appFile, maxSize, maxAge, maxBackups, compress)
	errLumberjackHook := lumberjackHook(errFile, maxSize, maxAge, maxBackups, compress)

	appZapCore := newZapCore(&appLumberjackHook, &appLogLevel)
	errZapCore := newZapCore(&errLumberjackHook, &errLogLevel)
	consoleCore := newZapCore(os.Stdout, &appLogLevel)

	var core zapcore.Core
	if showConsole {
		core = zapcore.NewTee(appZapCore, errZapCore, consoleCore)
	} else {
		core = zapcore.NewTee(appZapCore, errZapCore)
	}

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = &Logger{
		Zap:      zapLogger,
		ZapSugar: zapLogger.Sugar(),
		AppLevel: &appLogLevel,
	}

	config.RegisterConfigChangedFunc(onConfigChange)
}
