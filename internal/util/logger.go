package util

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger zerolog.Logger

func SetupLogger() {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack //nolint:reassign

	multi := zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		&lumberjack.Logger{
			Filename:   "/logs/airtime.log",
			MaxSize:    1,    // Max size in MB before log is rotated
			MaxBackups: 3,    // Max number of old log files to keep
			MaxAge:     28,   // Max age in days to retain log files
			Compress:   true, // Compress/zip old log files
		},
	)

	logger = zerolog.New(multi).With().Timestamp().Logger()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Fatal() *zerolog.Event {
	return logger.Fatal()
}

func Debug() *zerolog.Event {
	return logger.Debug()
}
