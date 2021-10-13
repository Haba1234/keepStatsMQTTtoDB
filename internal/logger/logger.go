package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger структура логгера.
type Logger struct {
	Log *logrus.Logger
}

// TODO добавить поддержку записи в файл.

// NewLogger конструктор.
func NewLogger() *Logger {
	log := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
	}
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableColors:    false,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})
	return &Logger{log}
}

func (l *Logger) Info(arg ...interface{}) {
	l.Log.Info(arg)
}

func (l *Logger) Infof(format string, arg ...interface{}) {
	l.Log.Infof(format, arg)
}

func (l *Logger) Error(arg ...interface{}) {
	l.Log.Error(arg)
}

func (l *Logger) Errorf(format string, arg ...interface{}) {
	l.Log.Errorf(format, arg)
}

func (l *Logger) Debug(arg ...interface{}) {
	l.Log.Debug(arg)
}

func (l *Logger) Debugf(format string, arg ...interface{}) {
	l.Log.Debugf(format, arg)
}
