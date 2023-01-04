package logging

import "go.uber.org/zap/zapcore"

const (
	termTimeFormat = "[01-02|15:04:05.000]"
)

var (
	defaultEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	termTimeEncoder = zapcore.TimeEncoderOfLayout(termTimeFormat)
)

func newTermEncoderConfig(lvEncoder zapcore.LevelEncoder) zapcore.EncoderConfig {
	cfg := defaultEncoderConfig
	cfg.EncodeLevel = lvEncoder
	cfg.EncodeTime = termTimeEncoder
	cfg.ConsoleSeparator = " "

	return cfg
}

func consoleColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := levelToColorString[Level(l)]
	if !ok {
		s = unknownLevelColor.Wrap(l.String())
	}
	enc.AppendString(s)
}

func ConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(newTermEncoderConfig(consoleColorLevelEncoder))
}
