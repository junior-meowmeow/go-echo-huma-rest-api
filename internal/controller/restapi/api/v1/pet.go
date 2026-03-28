package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterPetGroup(public huma.API, protected huma.API, h *handler.Handlers) {
	publicGroup := huma.NewGroup(public, "/pets")
	protectedGroup := huma.NewGroup(protected, "/pets")

	RegisterPetRoutes(publicGroup, protectedGroup, h)
}

func RegisterPetRoutes(public huma.API, protected huma.API, h *handler.Handlers) {
	huma.Register(public, huma.Operation{
		OperationID: "get-availble-pets",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "Get Available Pets",
		Description: "Get available pets.",
		Tags:        []string{"Pets"},
	}, h.Pet.GetAvailablePets)

	huma.Register(public, huma.Operation{
		OperationID: "get-pet-by-id",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Get Pet",
		Description: "Get a pet by ID.",
		Tags:        []string{"Pets"},
	}, h.Pet.GetPetByID)
}
