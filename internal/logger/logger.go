package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupLogger sets up the logger with the desired output format.
func SetupLogger() error {
	file, err := os.OpenFile("gofetch.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	log.Logger = log.Output(zerolog.SyncWriter(file))

	return nil
}

// SetLevel sets the global log level.
func SetLevel(levelStr string) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)
}
