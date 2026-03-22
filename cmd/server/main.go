package main

import (
	"context"
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
	server := http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.App.Port),
		Handler:     application.Router,
		BaseContext: func(net.Listener) context.Context { return notifyCtx },
	}

	// Run the server
	go func() {
		log.Printf("Starting server on port %d...\n", cfg.App.Port)
		log.Printf("API documentation is hosted at http://localhost:%d%s/docs\n", cfg.App.Port, cfg.App.APIBasePath)
		server.ListenAndServe()
	}()

	<-notifyCtx.Done()
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown
	server.Shutdown(shutdownCtx)
	application.GracefulShutdown(shutdownCtx)
	log.Println("Server exited gracefully.")
}
