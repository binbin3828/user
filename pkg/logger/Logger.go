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

var SugarLogger *zap.SugaredLogger

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	SugarLogger = logger.Sugar()
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
		Filename:   "../log/user.log", // 日志文件路径
		MaxSize:    1,                 // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 300,               // 日志文件最多保存多少个备份
		MaxAge:     30,                // 文件最多保存多少天
		Compress:   true,              // 是否压缩
		LocalTime:  true,              // 是否用当地时间命名文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
