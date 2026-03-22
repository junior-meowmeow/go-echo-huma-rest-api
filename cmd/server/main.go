package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/app"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"

	"github.com/danielgtaylor/huma/v2/humacli"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func main() {
	// Load configurations
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
	}

	// Create a CLI app which takes options
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {

		// Options overrides configurations
		if options.Port != 8888 {
			cfg.App.Port = options.Port
		}

		// Initialize Application
		application, err := app.NewApplication(context.Background(), cfg)
		if err != nil {
			log.Fatalf("Failed to initialize application: %v", err)
		}

		// Create a HTTP server
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.App.Port),
			Handler: application.Router,
		}

		hooks.OnStart(func() {
			log.Printf("Starting server on port %d...\n", cfg.App.Port)
			log.Printf("API documentation is hosted at http://localhost:%d%s/docs\n", cfg.App.Port, cfg.App.APIBasePath)
			server.ListenAndServe()
		})

		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			application.GracefulShutdown(ctx)
			log.Println("Server exited gracefully.")
		})
	})

	cli.Run()
}
