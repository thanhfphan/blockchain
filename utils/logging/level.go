package logging

import (
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

const (
	Verbo Level = iota - 9
	Debug
	Trace
	Info
	Warn
	Error
	Fatal
	Off
)

func (l Level) String() string {
	switch l {
	case Fatal:
		return "FATAL"
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	case Verbo:
		return "VERBO"
	case Off:
		return "OFF"
	default:
		return "UNKNOWN"
	}
}
