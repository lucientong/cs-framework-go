// Package logger provides to users directly
package logger

import "go.uber.org/zap"

// Debug zap原生格式日志
func Debug(msg string, fields ...zap.Field) {
	logger.Zap.Debug(msg, fields...)
}

// Debugf printf格式日志
func Debugf(template string, args ...interface{}) {
	logger.ZapSugar.Debugf(template, args)
}

// Info zap原生格式日志
func Info(msg string, fields ...zap.Field) {
	logger.Zap.Info(msg, fields...)
}

// Infof printf格式日志
func Infof(template string, args ...interface{}) {
	logger.ZapSugar.Infof(template, args)
}

// Warn zap原生格式日志
func Warn(msg string, fields ...zap.Field) {
	logger.Zap.Warn(msg, fields...)
}

// Warnf printf格式日志
func Warnf(template string, args ...interface{}) {
	logger.ZapSugar.Warnf(template, args)
}

// Error zap原生格式日志
func Error(msg string, fields ...zap.Field) {
	logger.Zap.Error(msg, fields...)
}

// Errorf printf格式日志
func Errorf(template string, args ...interface{}) {
	logger.ZapSugar.Errorf(template, args...)
}

// Panic zap原生格式日志
func Panic(msg string, fields ...zap.Field) {
	logger.Zap.Panic(msg, fields...)
}

// Panicf printf格式日志
func Panicf(template string, args ...interface{}) {
	logger.ZapSugar.Panicf(template, args...)
}

// Fatal zap原生格式日志
func Fatal(msg string, fields ...zap.Field) {
	logger.Zap.Fatal(msg, fields...)
}

// Fatalf printf格式日志
func Fatalf(template string, args ...interface{}) {
	logger.ZapSugar.Fatalf(template, args...)
}
