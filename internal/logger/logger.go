package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type AppLogger interface {
	Info(string)
	Error(string, error)
	Debug(string)
	Panic(error)
}

type appLogger struct {
	logger zerolog.Logger
}

func (l *appLogger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *appLogger) Error(message string, err error) {
	l.logger.Error().Err(err).Msg(message)
}

func (l *appLogger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *appLogger) Panic(err error) {
	l.Error("Panic error", err)
	l.logger.Panic().Msg(err.Error())
}

func NewAppLogger() AppLogger {
	// zerolog.TimeFieldFormat = time.RFC1123
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &appLogger{logger: log}
}
