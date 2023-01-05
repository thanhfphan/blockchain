package logging

import (
	"strings"

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

	FatalStr = "FATAL"
	ErrorStr = "ERROR"
	WarnStr  = "WARN"
	InfoStr  = "INFO"
	DebugStr = "DEBUG"
	VerboStr = "VERBO"
	OffStr   = "OFF"
)

func (l Level) String() string {
	switch l {
	case Fatal:
		return FatalStr
	case Error:
		return ErrorStr
	case Warn:
		return WarnStr
	case Info:
		return InfoStr
	case Debug:
		return DebugStr
	case Verbo:
		return VerboStr
	case Off:
		return OffStr
	default:
		return "UNKNOWN"
	}
}

func ToLevel(str string) Level {
	levelStr := strings.ToUpper(str)
	switch levelStr {
	case FatalStr:
		return Fatal
	case ErrorStr:
		return Error
	case WarnStr:
		return Warn
	case InfoStr:
		return Info
	case DebugStr:
		return Debug
	case VerboStr:
		return Verbo
	case OffStr:
		return Off
	default:
		return Off
	}
}
