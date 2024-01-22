package server

import (
	"GraphQL/configs"
	"GraphQL/internal/logger"
	"GraphQL/pkg/connection"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	cfg    *configs.Config
	pgDB   connection.DBops
	server *http.Server
}

func NewServer(cfg *configs.Config, db connection.DBops) *Server {
	return &Server{
		cfg:  cfg,
		pgDB: db,
	}
}

func (s *Server) Run(ctx context.Context) error {
	err := s.MapHandlers()
	if err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Infof(ctx, "Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := s.pgDB.Close(); err != nil {
			logger.Errorf(ctx, "Error closing Postgres database: %s", err)
		}

		if err := s.server.Shutdown(ctx); err != nil {
			logger.Errorf(ctx, "Server shutdown error: %v", err)
		}
	}()

	logger.Infof(ctx, "Server is listening on port %s", s.cfg.Server.Port)

	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	logger.Infof(ctx, "Server closed properly")

	return nil
}
