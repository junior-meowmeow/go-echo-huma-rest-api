package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterFileGroup(api huma.API, h *handler.Handlers) {
	fileGroup := huma.NewGroup(api, "/files")

	RegisterFileRoutes(fileGroup, h)
}

func RegisterFileRoutes(api huma.API, h *handler.Handlers) {
	huma.Register(api, huma.Operation{
		OperationID: "upload-file",
		Method:      http.MethodPost,
		Path:        "/upload",
		Summary:     "Upload file to file storage",
		Description: "Upload a file to file storage.",
		Tags:        []string{"Files"},
	}, h.File.UploadFile)

	huma.Register(api, huma.Operation{
		OperationID: "get-file-download-link",
		Method:      http.MethodGet,
		Path:        "/download/{id}",
		Summary:     "Get file from file storage",
		Description: "Get a file from file storage.",
		Tags:        []string{"Files"},
	}, h.File.GetFileDownloadLink)
}
