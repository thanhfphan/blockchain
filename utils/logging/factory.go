package logging

import (
	"fmt"
	"os"
	"sync"
)

var _ Factory = (*factory)(nil)

type Factory interface {
	Make(name string) (Logger, error)
	Close()
}

type factory struct {
	config Config
	lock   sync.RWMutex

	loggers map[string]Logger
}

func NewFactory(config Config) Factory {
	return &factory{
		config:  config,
		loggers: make(map[string]Logger),
	}
}

func (f *factory) makeLogger(cfg Config) (Logger, error) {
	if _, ok := f.loggers[cfg.LoggerName]; ok {
		return nil, fmt.Errorf("logger with name=%s already exists", cfg.LoggerName)
	}

	consoleEnc := ConsoleEncoder()
	consoleCore := NewWrappedCore(cfg.LogLevel, os.Stdout, consoleEnc)

	l := NewLogger(cfg.Prefix, consoleCore)
	f.loggers[cfg.LoggerName] = l

	return l, nil
}

func (f *factory) Make(name string) (Logger, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	config := f.config
	config.LoggerName = name
	return f.makeLogger(config)
}

func (f *factory) Close() {
	f.lock.Lock()
	defer f.lock.Unlock()

	for _, l := range f.loggers {
		l.Stop()
	}
}
