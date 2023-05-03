package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func CreateLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	runLogFile, _ := os.OpenFile(
		"dev.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, runLogFile)
	return zerolog.New(multi).With().Timestamp().Logger()
}
