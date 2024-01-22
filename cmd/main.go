package main

import (
	"GraphQL/configs"
	"GraphQL/internal/logger"
	"GraphQL/internal/server"
	"GraphQL/pkg/connection"
	"context"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	zapLogger, err := logger.InitLogger()
	if err != nil {
		log.Fatalf("Error initializing default logger %v", err)
	}

	logger.SetGlobal(zapLogger.With(zap.String("component", "main")))

	cfg, err := configs.LoadConfig()
	if err != nil {
		logger.Fatalf(ctx, "Could not import environment variables. %v", err)
	}

	dbManager, err := connection.NewDB(ctx, cfg)
	if err != nil {
		logger.Fatalf(ctx, "Could not connect database because of %v.", err)
	}

	app := server.NewServer(cfg, dbManager)

	if err := app.Run(ctx); err != nil {
		logger.Fatalf(ctx, "Cannot start server: %v", err)
	}
}
