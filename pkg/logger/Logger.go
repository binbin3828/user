/*
 * @Autor: Bobby
 * @Description: lib for logger
 * @Date: 2022-06-06 17:24:02
 * @LastEditTime: 2022-06-09 22:24:07
 * @FilePath: \user\pkg\logger\Logger.go
 */

package logger

import (
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 日志接口，用于依赖注入
type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return &ZapLogger{sugar: logger.Sugar()}
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Sync 刷新日志缓冲
func (l *ZapLogger) Sync() error {
	return l.sugar.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("[" + time.Format("2006-01-02 15:04:05") + "]")
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "../log/user.log",
		MaxSize:    1,
		MaxBackups: 300,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
