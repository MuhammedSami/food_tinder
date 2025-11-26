package app

import (
	"context"
	"errors"
	"fmt"
	"foodjiassignment/config"
	"foodjiassignment/internal/api"
	apiErr "foodjiassignment/internal/api/errors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	API    *api.APIInterface
	DB     *gorm.DB
	Config *config.Config
}

func NewApp(db *gorm.DB, cfg *config.Config) *App {
	app := &App{
		DB:     db,
		Config: cfg,
	}

	redisClient := app.GetRedisClient(cfg.Redis)

	app.API = api.NewApi(
		app.GetTinderFoodManager(),
		apiErr.NewApiError(),
		redisClient,
	)

	return app
}

func (a *App) ExposeWithGracefulShutDown(ctx context.Context) error {
	r := a.API.RegisterHandlers()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.Config.Api.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf("server is running on http://localhost%s", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("server error: %v", err)
	case sig := <-stop:
		log.Printf("received shutdown signal: %v", sig)
	}

	log.Println("server is shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %+v", err)
	}

	// shutdown DB as well, no need for hanging connection on network layer
	sqlDB, err := a.DB.DB()
	if err != nil {
		log.Printf("failed to get sql.DB: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}

	log.Println("server exited properly")

	return nil
}
