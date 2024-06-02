package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupLogger sets up the logger with the desired output format.
func SetupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	})
}

// SetLevel sets the global log level.
func SetLevel(levelStr string) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)
}
