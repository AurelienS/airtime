package log

import (
	"os"

	"github.com/rs/zerolog" // Global logger
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupLogger initializes and returns a global logger
func SetupLogger() zerolog.Logger {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	multi := zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		&lumberjack.Logger{
			Filename:   "./logs/app.log",
			MaxSize:    1,    // Max size in MB before log is rotated
			MaxBackups: 3,    // Max number of old log files to keep
			MaxAge:     28,   // Max age in days to retain log files
			Compress:   true, // Compress/zip old log files
		},
	)

	// Return the configured global logger
	return zerolog.New(multi).With().Timestamp().Logger()
}
