// package logger provides a zerolog.Logger instance.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New creates a new logger instance.
func New() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	// Set logger output to console if not in release mode
	if os.Getenv("GIN_MODE") != "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
	return log.Logger
}
