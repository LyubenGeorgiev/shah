package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LyubenGeorgiev/shah/cache"
	"github.com/LyubenGeorgiev/shah/chess"
	"github.com/LyubenGeorgiev/shah/db"
	"github.com/gorilla/mux"
)

type App struct {
	router  *mux.Router
	Storage db.Storage
	Cache   cache.Cache
	Manager *chess.Manager
}

func New() *App {
	app := &App{
		router:  mux.NewRouter().StrictSlash(true),
		Storage: db.NewPostgresStorage(),
		Cache:   cache.NewRedisCache(),
	}
	app.Manager = chess.NewManager(app.Storage)

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           a.router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := a.Cache.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("Failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.Cache.Close(); err != nil {
			fmt.Println("Failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			ch <- fmt.Errorf("Failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
