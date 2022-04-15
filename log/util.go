package log

import (
	"cs/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"runtime"
	"strings"
	"time"
)

var logLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLogLevel(level string) zapcore.Level {
	if l, ok := logLevelMap[level]; ok {
		return l
	}
	return zapcore.InfoLevel
}

// lumberjackHook 日志切割
func lumberjackHook(fileName string, maxSize, maxAge, maxBackups int, compress bool) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		Compress:   compress,
	}
}

// isO8601TimeEncoder 修改时间格式
func isO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
}

// capitalLevelEncoder 修改日志等级格式
func capitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + l.CapitalString() + "]")
}

// callerEncoder 添加日志调用者函数名
func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	funcPath := runtime.FuncForPC(caller.PC).Name()
	funcPaths := strings.Split(funcPath, "/")
	funcName := funcPaths[len(funcPaths)-1:][0]
	enc.AppendString("[" + caller.TrimmedPath() + " " + funcName + "]")
}

func newZapCore(w io.Writer, atomicLevel *zap.AtomicLevel) zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    capitalLevelEncoder,
		EncodeTime:     isO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   callerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	writer := zapcore.AddSync(w)
	return zapcore.NewCore(encoder, writer, atomicLevel)
}

// onConfigChange 配置文件更改后，热变更日志级别
func onConfigChange(name string, op uint32) {
	level := config.Use("log").MustString("level", "info")
	logger.SetAppLogLevel(level)
}
