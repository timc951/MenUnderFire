package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

// Config holds logger configuration
type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	TimeFormat string // RFC3339, Unix, etc.
}

// Init initializes the global logger with the given configuration
func Init(cfg Config) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set time format
	if cfg.TimeFormat == "" {
		zerolog.TimeFieldFormat = time.RFC3339
	} else {
		zerolog.TimeFieldFormat = cfg.TimeFormat
	}

	// Set output format
	var output io.Writer = os.Stdout
	if cfg.Format == "console" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	Log = zerolog.New(output).With().Timestamp().Caller().Logger()
}

// InitDefault initializes the logger with sensible defaults for development
func InitDefault() {
	Init(Config{
		Level:  "debug",
		Format: "console",
	})
}

// InitProduction initializes the logger for production (JSON format)
func InitProduction() {
	Init(Config{
		Level:  "info",
		Format: "json",
	})
}

// Helper functions for common logging patterns

func Debug() *zerolog.Event {
	return Log.Debug()
}

func Info() *zerolog.Event {
	return Log.Info()
}

func Warn() *zerolog.Event {
	return Log.Warn()
}

func Error() *zerolog.Event {
	return Log.Error()
}

func Fatal() *zerolog.Event {
	return Log.Fatal()
}

// WithComponent returns a logger with a component field
func WithComponent(component string) zerolog.Logger {
	return Log.With().Str("component", component).Logger()
}
