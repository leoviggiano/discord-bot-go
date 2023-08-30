package log

import (
	"bot_test/config"

	"github.com/sirupsen/logrus"
)

type Option func(logger *logrus.Logger)

func WithTextFormatter() Option {
	return func(logger *logrus.Logger) {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:            true,
			DisableLevelTruncation: true,
			DisableTimestamp:       true,
		})
	}
}

func NewLogger(level logrus.Level, options ...Option) logrus.FieldLogger {
	logger := logrus.New()
	logger.Level = level

	for _, option := range options {
		option(logger)
	}

	if config.Environment() == config.Production {
		logger.Level = logrus.ErrorLevel
	}

	logger.Info("[Log] Started with success")
	return logger
}

func IfError(logger logrus.FieldLogger, err error) {
	if err != nil {
		logger.Error(err)
	}
}
