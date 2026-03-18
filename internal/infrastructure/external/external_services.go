package external

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external/petstore"
	petStoreClient "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external/petstore/client"
)

type ExternalServices struct {
	PetService port.PetService
}

func NewExternalServices(petStoreClient *petStoreClient.ClientWithResponses) *ExternalServices {
	return &ExternalServices{
		PetService: petstore.NewPetService(petStoreClient),
	}
}
