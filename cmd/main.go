package main

import (
	"context"
	"foodjiassignment/config"
	"foodjiassignment/internal/app"
	"foodjiassignment/internal/storage"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("failed to validate config %-+v", err)
	}

	ctx := context.Background()

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(db, cfg)

	if err := a.ExposeWithGracefulShutDown(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
