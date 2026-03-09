package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/db"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/handlers"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"

	"github.com/danielgtaylor/huma/v2/humacli"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func main() {
	// Load config from environment variables
	cfg := config.NewConfig()

	// Initialize DB Clients
	mongoClient, err := db.NewMongoDBClient(cfg.MongoUser, cfg.MongoPass, cfg.MongoHost, cfg.MongoPort)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.DisconnectMongoDB(mongoClient)

	s3Client, err := db.NewS3Client(cfg.S3Endpoint)
	if err != nil {
		log.Fatalf("Failed to connect to S3: %v", err)
	}

	// Initialize Repositories
	mongoDB := mongoClient.Database(cfg.DBName)

	repositories := repositories.NewRepositories(mongoDB, s3Client, cfg.S3Bucket)

	// Initialize Handler
	h := handlers.NewHandler(repositories)

	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {

		// Create a new router and register APIs (from internal/api)
		router := api.NewRouter(h, cfg.APIBasePath)

		port := cfg.Port
		if options.Port != 8888 { // CLI flag overrides env
			port = options.Port
		}

		// Create a HTTP server.
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}

		hooks.OnStart(func() {
			log.Printf("Starting server on port %d...\n", port)
			log.Printf("API documentation is hosted at http://localhost:%d%s/docs\n", port, cfg.APIBasePath)
			server.ListenAndServe()
		})

		// Tell the CLI how to stop your server.
		hooks.OnStop(func() {
			// Give the server 5 seconds to gracefully shut down, then give up.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	cli.Run()
}
