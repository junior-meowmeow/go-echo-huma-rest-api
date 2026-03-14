package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type PetHandler interface {
	GetAvailablePets(ctx context.Context, request *schema.GetAvailablePetsRequest) (*schema.GetAvailablePetsResponse, error)
	GetPetByID(ctx context.Context, request *schema.GetPetByIDRequest) (*schema.GetPetByIDResponse, error)
}

type petHandler struct {
	PetUseCase usecase.PetUseCase
}

func NewPetHandler(petUseCase usecase.PetUseCase) *petHandler {
	return &petHandler{
		PetUseCase: petUseCase,
	}
}

func (h *petHandler) GetAvailablePets(ctx context.Context, request *schema.GetAvailablePetsRequest) (*schema.GetAvailablePetsResponse, error) {
	pets, err := h.PetUseCase.GetAvailablePets(ctx)
	if err != nil {
		return nil, err
	}

	resp := schema.GetAvailablePetsResponse{}
	resp.Body.Data = convertPets(pets)

	return &resp, nil
}

func (h *petHandler) GetPetByID(ctx context.Context, request *schema.GetPetByIDRequest) (*schema.GetPetByIDResponse, error) {
	pet, err := h.PetUseCase.GetPetByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	resp := schema.GetPetByIDResponse{}
	resp.Body = convertPet(pet)

	return &resp, nil
}

func convertPets(pets []entity.Pet) []schema.Pet {
	petOutputs := make([]schema.Pet, len(pets))
	for i, p := range pets {
		petOutputs[i] = convertPet(p)
	}
	return petOutputs
}

func convertPet(pet entity.Pet) schema.Pet {
	return schema.Pet{
		ID:   pet.ID,
		Name: pet.Name,
		Category: schema.PetCategory{
			ID:   pet.Category.ID,
			Name: pet.Category.Name,
		},
		PhotoURLs: pet.PhotoURLs,
		Status:    string(pet.Status),
		Tags:      pet.Tags,
	}
}
