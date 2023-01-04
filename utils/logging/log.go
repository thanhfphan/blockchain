package logging

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*log)(nil)

type log struct {
	internalLogger *zap.Logger
	wrappedCore    []WrappedCore
}

type WrappedCore struct {
	Core           zapcore.Core
	Writer         io.WriteCloser
	WriterDisabled bool
	AtomicLevel    zap.AtomicLevel
}

func NewWrappedCore(level Level, rw io.WriteCloser, encoder zapcore.Encoder) WrappedCore {
	atomicLevel := zap.NewAtomicLevelAt(zapcore.Level(level))

	core := zapcore.NewCore(encoder, zapcore.AddSync(rw), atomicLevel)
	return WrappedCore{AtomicLevel: atomicLevel, Core: core, Writer: rw}
}

func newZapLogger(prefix string, wrappedCores ...WrappedCore) *zap.Logger {
	cores := make([]zapcore.Core, len(wrappedCores))
	for i, wc := range wrappedCores {
		cores[i] = wc.Core
	}
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
	if prefix != "" {
		logger = logger.Named(prefix)
	}

	return logger
}

func NewLogger(prefix string, wrappedCores ...WrappedCore) Logger {
	return &log{
		internalLogger: newZapLogger(prefix, wrappedCores...),
		wrappedCore:    wrappedCores,
	}
}

func (l *log) Stop() {
	for _, wc := range l.wrappedCore {
		_ = wc.Writer.Close()
	}
}

func (l *log) log(level Level, msg string) {
	if ce := l.internalLogger.Check(zapcore.Level(level), msg); ce != nil {
		ce.Write()
	}
}

func (l *log) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(Fatal, msg)
}

func (l *log) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(Error, msg)
}

func (l *log) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(Warn, msg)

}

func (l *log) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(Info, msg)
}

func (l *log) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(Debug, msg)
}

func (l *log) Verbof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(Verbo, msg)
}
