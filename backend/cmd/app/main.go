package main

import (
	"context"
	"errors"
	"golbugames/config"
	"golbugames/internal/api/middleware"
	"golbugames/internal/api/router"
	"golbugames/internal/database"
	"golbugames/internal/websocket"
	"golbugames/internal/websocket/multiplayer"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	settings, err := config.Load()
	if err != nil {
		return err
	}
	if err := websocket.ConfigureMaxPayloadSize(settings.Payload.WebSocketMaxBytes); err != nil {
		return err
	}

	jwtManager, err := middleware.NewJWTManager(settings.JWT)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = database.InitDB(ctx, database.PoolSettings{
		URL:                  settings.Database.URL,
		MinConnections:       settings.Database.MinConnections,
		MaxConnections:       settings.Database.MaxConnections,
		MaxConnectionAge:     settings.Database.MaxConnectionAge,
		MaxConnectionIdle:    settings.Database.MaxConnectionIdle,
		ConnectTimeout:       settings.Database.ConnectTimeout,
		QueryTimeout:         settings.Database.QueryTimeout,
		BulkOperationTimeout: settings.Database.BulkOperationTimeout,
	})
	if err != nil {
		return err
	}
	defer database.CloseDB()

	err = config.InitGridGeneration(ctx)
	if err != nil {
		return err
	}

	hubManager := multiplayer.NewHubManager()
	go hubManager.MatchmakingLoop(ctx)

	r := router.NewRouter(hubManager, jwtManager)
	r.InitRoutes()

	server := &http.Server{
		Addr:              settings.HTTP.Address,
		Handler:           middleware.NewHTTPHandler(r, settings),
		ReadTimeout:       settings.HTTP.ReadTimeout,
		ReadHeaderTimeout: settings.HTTP.ReadHeaderTimeout,
		WriteTimeout:      settings.HTTP.WriteTimeout,
		IdleTimeout:       settings.HTTP.IdleTimeout,
		MaxHeaderBytes:    settings.HTTP.MaxHeaderBytes,
	}

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), settings.HTTP.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
		}
	}()

	slog.Info("HTTP server listening", "address", settings.HTTP.Address)
	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-shutdownDone
	return nil
}
