package logger

import (
	"fmt"
	"io"
	"sync"
)

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
)

func New(out io.Writer, err io.Writer, level string) *Logger {
	var actualLevel Level
	if level == "" {
		actualLevel = InfoLevel
	} else {
		actualLevel = Level(level)
	}

	return &Logger{
		out:   out,
		err:   err,
		level: actualLevel,
		mutex: sync.Mutex{},
	}
}

type Logger struct {
	out   io.Writer
	err   io.Writer
	level Level
	mutex sync.Mutex
}

func (l *Logger) SetLevel(level Level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

func (l *Logger) Debug(msg string, params ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.level > DebugLevel {
		return
	}
	fmt.Fprintf(l.out, "%s: %s\n", DebugLevel, fmt.Sprintf(msg, params...)) //nolint:errcheck
}

func (l *Logger) Info(msg string, params ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.level > InfoLevel {
		return
	}
	fmt.Fprintf(l.out, "%s: %s\n", InfoLevel, fmt.Sprintf(msg, params...)) //nolint:errcheck
}

func (l *Logger) Warn(msg string, params ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.level > WarnLevel {
		return
	}
	fmt.Fprintf(l.out, "%s: %s\n", WarnLevel, fmt.Sprintf(msg, params...)) //nolint:errcheck
}

func (l *Logger) Error(msg string, params ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.level > ErrorLevel {
		return
	}
	fmt.Fprintf(l.err, "%s: %s\n", ErrorLevel, fmt.Sprintf(msg, params...)) //nolint:errcheck
}
