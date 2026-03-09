package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/api"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/handlers"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories"
)

func main() {
	// Initialize Handlers without Repositories
	handlers := handlers.NewHandlers(&repositories.Repositories{})

	// Initialize Router
	router := api.NewRouter(handlers, "")

	// Request API Documentations and Write them to docs/
	requestAndWriteDocs(router)
}

func requestAndWriteDocs(router http.Handler) {
	fmt.Println("Generating API Documentations...")

	// Ensure docs directory exists
	if err := os.MkdirAll("docs", 0755); err != nil {
		fmt.Printf("❌ Failed to create docs directory: %v\n", err)
		os.Exit(1)
	}

	docsToGen := map[string]string{
		"/openapi.json":     "docs/openapi.json",
		"/openapi.yaml":     "docs/openapi.yaml",
		"/openapi-3.0.json": "docs/openapi-3.0.json",
		"/openapi-3.0.yaml": "docs/openapi-3.0.yaml",
	}

	for path, filename := range docsToGen {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()

		// Request the documentation from the router
		router.ServeHTTP(w, req)

		// Check if the request was successful
		if w.Code != http.StatusOK {
			fmt.Printf("⚠️ Failed to request %s (Status: %d)\n", filename, w.Code)
			continue
		}

		outBytes := w.Body.Bytes()

		// Pretty print JSON file
		if strings.HasSuffix(filename, ".json") {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, outBytes, "", "  "); err == nil {
				outBytes = prettyJSON.Bytes()
			} else {
				fmt.Printf("⚠️ Failed to format JSON for %s: %v\n", filename, err)
			}
		}

		// Write the documentation file
		if err := os.WriteFile(filename, outBytes, 0644); err != nil {
			fmt.Printf("❌ Failed to write %s: %v\n", filename, err)
			os.Exit(1)
		}

		fmt.Printf("✅ Generated %s\n", filename)
	}

	fmt.Println("Generated All API Documentations.")
}
