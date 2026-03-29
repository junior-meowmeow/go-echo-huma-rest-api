package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		log.Fatalf("Failed to load configurations: %v", err)
	}

	// Initialize Application
	application, err := app.NewApplication(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
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
		log.Printf("Starting server on port %d...\n", cfg.App.Port)
		log.Printf("API documentation is hosted at http://localhost:%d%s/docs\n", cfg.App.Port, cfg.App.APIBasePath)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start or crashed: %v\n", err)
		}
	}()

	<-notifyCtx.Done()
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	const shutdownTimeout = 5 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown with error: %v\n", err)
	}
	application.GracefulShutdown(shutdownCtx)
	log.Println("Server exited gracefully.")
}
