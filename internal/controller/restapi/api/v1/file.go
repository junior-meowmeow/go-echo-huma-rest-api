package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterFileGroup(public huma.API, protected huma.API, h *handler.Handlers) {
	publicGroup := huma.NewGroup(public, "/files")
	protectedGroup := huma.NewGroup(protected, "/files")

	RegisterFileRoutes(publicGroup, protectedGroup, h)
}

//revive:disable-next-line:unused-parameter // Keeps a consistent signature across all route registration functions.
func RegisterFileRoutes(public huma.API, protected huma.API, h *handler.Handlers) {
	huma.Register(public, huma.Operation{
		OperationID: "upload-file",
		Method:      http.MethodPost,
		Path:        "/upload",
		Summary:     "Upload file to file storage",
		Description: "Upload a file to file storage.",
		Tags:        []string{"Files"},
	}, h.File.UploadFile)

	huma.Register(public, huma.Operation{
		OperationID: "get-file-download-link",
		Method:      http.MethodGet,
		Path:        "/download/{id}",
		Summary:     "Get file from file storage",
		Description: "Get a file from file storage.",
		Tags:        []string{"Files"},
	}, h.File.GetFileDownloadLink)
}
