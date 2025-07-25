package klog

import (
	"flag"

	"go.uber.org/zap"
	"k8s.io/klog"
)

type infoLogger struct {
	logger *zap.Logger
}

type warnLogger struct {
	logger *zap.Logger
}

type errorLogger struct {
	logger *zap.Logger
}

type fatalLogger struct {
	logger *zap.Logger
}

func InitLogger(zapLogger *zap.Logger) {
	fs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(fs)
	defer klog.Flush()

	klog.SetOutputBySeverity("INFO", &infoLogger{logger: zapLogger})
	klog.SetOutputBySeverity("WARNING", &infoLogger{logger: zapLogger})
	klog.SetOutputBySeverity("FATAL", &infoLogger{logger: zapLogger})
	klog.SetOutputBySeverity("ERROR", &infoLogger{logger: zapLogger})

	_ = fs.Set("skip_header", "true")
	_ = fs.Set("logtostderr", "false")
}

func (l *infoLogger) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p[:len(p)-1]))

	return len(p), nil
}

func (l *warnLogger) Write(p []byte) (n int, err error) {
	l.logger.Warn(string(p[:len(p)-1]))

	return len(p), nil
}

func (l *errorLogger) Write(p []byte) (n int, err error) {
	l.logger.Error(string(p[:len(p)-1]))

	return len(p), nil
}

func (l *fatalLogger) Write(p []byte) (n int, err error) {
	l.logger.Fatal(string(p[:len(p)-1]))

	return len(p), nil
}
