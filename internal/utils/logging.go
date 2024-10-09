package utils

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once sync.Once

func ConfigureLoggingFromEnvOnce() {
	once.Do(func() {
		logLevel := GetEnv("LOG_LEVEL", "debug")
		logFormat := GetEnv("LOG_FORMAT", "json")

		ConfigureLogging(logLevel, logFormat)
	})
}

func ConfigureLogging(logLevel string, logFormat string) {
	log.Logger = log.With().Timestamp().Caller().Logger().Level(zerolog.DebugLevel)

	if logLevel != "" {
		level, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Panic().Err(err).Str("level", logLevel).Msg("Invalid log level")
		}
		log.Logger = log.Logger.Level(level)
	}

	if logFormat == "json" {
		// Defaults to JSON already nothing to do
	} else if logFormat == "" || logFormat == "pretty" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    false,
			TimeFormat: time.RFC3339,
		})
	} else {
		log.Panic().Str("log_format", logFormat).Msg("Unknown log format")
	}
}
