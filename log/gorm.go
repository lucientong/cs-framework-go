// Package log 整合gorm日志到zap
package log

import (
	"context"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	slowThreshold time.Duration // 慢查询阈值
	logger        *Logger
}

// NewGormLogger 基于zap的gorm日志
func NewGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		slowThreshold: slowThreshold,
		logger:        logger,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(gormLogLevel gormLogger.LogLevel) gormLogger.Interface {
	return l
}

// Info 打印Info日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.ZapSugar.Infof(msg, data...)
}

// Warn 打印Warn日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.ZapSugar.Warnf(msg, data...)
}

// Error 打印Info日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.ZapSugar.Errorf(msg, data...)
}

// Trace trace
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// todo
}
