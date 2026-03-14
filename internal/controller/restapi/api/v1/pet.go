package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterPetGroup(api huma.API, h *handler.Handlers) {
	bookGroup := huma.NewGroup(api, "/pets")

	RegisterPetRoutes(bookGroup, h)
}

func RegisterPetRoutes(api huma.API, h *handler.Handlers) {
	huma.Register(api, huma.Operation{
		OperationID: "get-availble-pets",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Get Available Pets",
		Description: "Get available pets.",
		Tags:        []string{"Pets"},
	}, h.Pet.GetAvailablePets)

	huma.Register(api, huma.Operation{
		OperationID: "get-pet-by-id",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Get Pet",
		Description: "Get a pet by ID.",
		Tags:        []string{"Pets"},
	}, h.Pet.GetPetByID)
}
