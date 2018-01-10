package log

import (
	"context"
	"os"
	"time"

	"github.com/marove2000/hack-and-pay/ctxutil"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

const (
	correlationId = "correlationId"
)

type Logger struct {
	ll     *logrus.Logger
	fields logrus.Fields
}

func init() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		TimestampFormat: time.RFC1123,
		FullTimestamp:   true,
		ForceFormatting: true,
	})
	logrus.SetOutput(os.Stdout)
}

func New(pkg string) *Logger {
	logger := &Logger{
		ll: logrus.StandardLogger(),
		fields: logrus.Fields{
			"pkg": pkg,
		},
	}

	return logger
}

func (l *Logger) ForFunc(ctx context.Context, fn string) *Logger {
	fields := logrus.Fields{
		"fn": fn,
	}

	corr := ctxutil.GetCorrelationId(ctx)
	if corr != "" {
		fields[correlationId] = corr
	}

	for k, v := range l.fields {
		fields[k] = v
	}

	return &Logger{
		ll:     l.ll,
		fields: fields,
	}
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	e := l.ll.WithField(key, value)

	for k, v := range l.fields {
		e.Data[k] = v
	}

	return e
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	e := l.ll.WithFields(fields)

	for k, v := range l.fields {
		e.Data[k] = v
	}

	return e
}

func (l *Logger) WithError(err error) *logrus.Entry {
	e := l.ll.WithError(err)

	for k, v := range l.fields {
		e.Data[k] = v
	}

	return e
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Infof(format, args...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Printf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Warnf(format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Warningf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Fatalf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.ll.WithFields(l.fields).Panicf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.ll.WithFields(l.fields).Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.ll.WithFields(l.fields).Info(args...)
}

func (l *Logger) Print(args ...interface{}) {
	l.ll.WithFields(l.fields).Print(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.ll.WithFields(l.fields).Warn(args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.ll.WithFields(l.fields).Warning(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.ll.WithFields(l.fields).Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.ll.WithFields(l.fields).Fatal(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.ll.WithFields(l.fields).Panic(args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.ll.WithFields(l.fields).Debugln(args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.ll.WithFields(l.fields).Infoln(args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.ll.WithFields(l.fields).Println(args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.ll.WithFields(l.fields).Warnln(args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.ll.WithFields(l.fields).Warningln(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.ll.WithFields(l.fields).Errorln(args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.ll.WithFields(l.fields).Fatalln(args...)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.ll.WithFields(l.fields).Panicln(args...)
}
