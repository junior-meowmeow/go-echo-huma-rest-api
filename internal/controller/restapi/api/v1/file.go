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
		OperationID: "upload-s3-file",
		Method:      http.MethodPost,
		Path:        "/upload",
		Summary:     "Upload file to object storage",
		Description: "Upload a file to object storage.",
		Tags:        []string{"Files"},
	}, h.File.UploadFile)

	huma.Register(api, huma.Operation{
		OperationID: "get-s3-file",
		Method:      http.MethodGet,
		Path:        "/download/{id}",
		Summary:     "Get file from object storage",
		Description: "Get a file from object storage.",
		Tags:        []string{"Files"},
	}, h.File.GetFileDownloadLink)
}
