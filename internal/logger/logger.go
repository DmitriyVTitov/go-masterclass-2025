package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(level zerolog.Level) {
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    true,
			TimeFormat: "2006-01-02 15:04:05",
		},
	).With().Timestamp().Logger().With().Caller().Logger()

	zerolog.SetGlobalLevel(level)
}
