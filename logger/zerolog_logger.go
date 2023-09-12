package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type zerologLogger struct {
	logger *zerolog.Logger
}

type zerologLoggerConfig struct {
	PrettyLogging bool
}

func NewZerologLoggerConfig() *zerologLoggerConfig {
	config := zerologLoggerConfig{}
	config.PrettyLogging = true

	return &config
}

func NewZerologLogger(config *zerologLoggerConfig) Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if config.PrettyLogging {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return &zerologLogger{logger: &logger}
}

func (zl *zerologLogger) Info(msg string) {
	zl.logger.Info().Msg(msg)
}

func (zl *zerologLogger) Infof(format string, args ...interface{}) {
	zl.logger.Info().Msgf(format, args...)
}

func (zl *zerologLogger) Error(msg string) {
	zl.logger.Error().Msg(msg)
}

func (zl *zerologLogger) Errorf(format string, args ...interface{}) {
	zl.logger.Error().Msgf(format, args...)
}
