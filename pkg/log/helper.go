package log

import (
	"context"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// convert key-value into zap fields
func handleFileds(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		return additional
	}

	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed zap field passed to logr", Any("zap filed", args[i]))

			break
		}

		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", Any("ignored key", args[i]))

			break
		}

		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			l.DPanic("non-string key argument passed to logging, ignoring all later arguments", Any("invalid key", key))
			break
		}
		fields = append(fields, Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = New(opts)
}

func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

func StdErrLogger() *log.Logger {
	if std == nil {
		return nil
	}

	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}

	return nil
}

func StdInfoLogger() *log.Logger {
	if std == nil {
		return nil
	}

	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

func V(level Level) InfoLogger {
	return std.V(level)
}

func WithValues(keysAndValues ...interface{}) Logger {
	return std.WithValues(keysAndValues...)
}

func WithName(s string) Logger {
	return std.WithName(s)
}

func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(Logger)
		}
	}

	return WithName("Unknown-Context")
}

func Flush() {
	std.Flush()
}

func ZapLogger() *zap.Logger {
	return std.zapLogger
}

func CheckIntLevel(level int32) bool {
	var zapLevel zapcore.Level

	if zapLevel < 5 {
		zapLevel = zapcore.InfoLevel
	} else {
		zapLevel = zapcore.DebugLevel
	}

	checkEntry := std.zapLogger.Check(zapLevel, "")
	return checkEntry != nil
}

func Debug(msg string, fields ...Field) {
	std.zapLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func Info(msg string, fields ...Field) {
	std.zapLogger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	std.zapLogger.Sugar().Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func Warn(msg string, fields ...Field) {
	std.zapLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...Field) {
	std.zapLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func Panic(msg string, fields ...Field) {
	std.zapLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...Field) {
	std.zapLogger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func L(ctx context.Context) *zapLogger {
	return std.L(ctx)
}
