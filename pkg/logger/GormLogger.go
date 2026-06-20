package logger

import (
	"context"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger        Logger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		l.Logger.Errorf("[%.2fms] [rows:%d] %s | error: %v", float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
	} else {
		l.Logger.Infof("[%.2fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
