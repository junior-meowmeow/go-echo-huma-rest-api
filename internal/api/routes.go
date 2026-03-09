package api

import (
	"net/http"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, h *handlers.Handlers) {
	// GET /health
	huma.Register(api, huma.Operation{
		OperationID:   "health-check",
		Method:        http.MethodGet,
		Path:          "/health",
		Summary:       "Get health status",
		Description:   "Get a health check status of the server.",
		Tags:          []string{"Monitoring"},
		DefaultStatus: 200,
	}, h.Health.GetHealthStatus)

	// GET /greeting/{name}
	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Get a greeting",
		Description: "Get a greeting for a person by name.",
		Tags:        []string{"Miscellaneous"},
	}, h.Greeting.GetGreeting)

	// GET /test-s3
	huma.Register(api, huma.Operation{
		OperationID: "list-s3-files",
		Method:      http.MethodGet,
		Path:        "/test-s3",
		Summary:     "Test S3 Connection",
		Description: "Directly lists files from the S3 bucket to verify connectivity.",
		Tags:        []string{"Miscellaneous"},
	}, h.Files.ListS3Files)
}
