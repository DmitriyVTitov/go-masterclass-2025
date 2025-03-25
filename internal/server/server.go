package server

import (
	"context"

	"ugc/internal/api"
	"ugc/internal/config"
	"ugc/internal/db"
	"ugc/internal/db/memdb"
	"ugc/internal/reviews"

	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg     *config.Config
	db      db.DB
	reviews *reviews.Service
	api     *api.API
}

func New(cfg *config.Config) (*Server, error) {
	s := Server{
		cfg: cfg,
	}

	s.db = memdb.New()

	s.reviews = reviews.New(s.db)

	s.api = api.New(*s.cfg, s.reviews)

	return &s, nil
}

func (s *Server) Run(ctx context.Context) error {
	// Init Telemetry SDK.
	shutdown, err := setupOTelSDK(ctx, s.cfg.TelemetryEndpoint)
	if err != nil {
		return err
	}
	defer func() {
		err = shutdown(ctx)
		if err != nil {
			log.Err(err).Msg("cannot shutdown OTel")
		}
	}()

	// Start API HTTP Server.
	return s.api.Serve(ctx)
}

func (s *Server) Shutdown() {
	log.Info().Msg("graceful server shutdown")
}
