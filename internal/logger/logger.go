package logger

import (
	"fmt"
	"os"

	"github.com/Haba1234/keepStatsMQTTtoDB/internal/config"
	"github.com/sirupsen/logrus"
)

// Log структура логгера.
type Log struct {
	Logger *logrus.Logger
}

// TODO добавить поддержку записи в файл.

// NewLogger конструктор.
func NewLogger(cfg config.LogConf) (*Log, error) {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableColors:    false,
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	}

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("logger. Error in settings (level: %s): %w", cfg.Level, err)
	}
	log.SetLevel(level)
	log.Debug("set level: ", level)

	return &Log{log}, nil
}

func (l *Log) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *Log) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *Log) Debug(args ...interface{}) {
	if len(args) >= 3 {
		switch args[0].(type) {
		case string:
			l.Logger.WithFields(logrus.Fields{
				fmt.Sprintf("%v", args[0]): args[1],
			}).Debug(args[2:])
			return
		default:
		}
	}
	l.Logger.Debug(args...)
}

func (l *Log) GetLevel() string {
	return l.Logger.Level.String()
}
