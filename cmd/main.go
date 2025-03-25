package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"ugc/internal/config"
	"ugc/internal/logger"
	"ugc/internal/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

func main() {
	logger.Init(zerolog.DebugLevel)
	configPath := pflag.StringP("config", "c", "", "path to config file")
	pflag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("can not load config")
	}

	s, err := server.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("can not create server")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	go func() {
		err = s.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-ctx.Done()
	s.Shutdown()
}
