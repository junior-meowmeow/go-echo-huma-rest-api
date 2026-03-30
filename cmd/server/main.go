package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/app"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
)

func main() {
	// Load configurations
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to load configurations", slog.Any("error", err))
		os.Exit(1)
	}

	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(cfg.Log.Level)); err != nil {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting application...")

	// Initialize Application
	application, err := app.NewApplication(context.Background(), cfg)
	if err != nil {
		slog.Error("Failed to initialize application", slog.Any("error", err))
		os.Exit(1)
	}

	notifyCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create a HTTP server
	const readHeaderTimeout = 5 * time.Second
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           application.Router,
		BaseContext:       func(net.Listener) context.Context { return notifyCtx },
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Run the server
	go func() {
		slog.Info(fmt.Sprintf("Starting server on port %d...", cfg.App.Port))
		slog.Debug(fmt.Sprintf("API documentation is hosted at http://localhost:%d%s/docs ", cfg.App.Port, cfg.App.APIBasePath))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start or crashed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-notifyCtx.Done()
	slog.Info("Shutdown signal received, initiating graceful shutdown...")

	const shutdownTimeout = 5 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Warn("Server shutdown with error", slog.Any("error", err))
	}
	application.GracefulShutdown(shutdownCtx)
	slog.Info("Server exited gracefully.")
}
