package petstore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external/petstore/client"
)

type petService struct {
	petStoreClient *client.ClientWithResponses
}

func NewPetService(petStoreClient *client.ClientWithResponses) *petService {
	return &petService{
		petStoreClient: petStoreClient,
	}
}

func (s *petService) GetPetByID(ctx context.Context, id int64) (entity.Pet, error) {
	resp, err := s.petStoreClient.GetPetByIdWithResponse(ctx, id)
	if err != nil {
		return entity.Pet{}, fmt.Errorf("network error calling petstore api: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNotFound {
			return entity.Pet{}, fmt.Errorf("pet not found")
		}
		return entity.Pet{}, fmt.Errorf("unexpected status code from petstore client: %d", resp.StatusCode())
	}

	if resp.JSON200 == nil {
		return entity.Pet{}, fmt.Errorf("received 200 OK but body was empty")
	}

	pet := mapClientPetToEntity(resp.JSON200)

	return pet, nil
}

func (s *petService) GetPetsByStatus(ctx context.Context, status entity.PetStatus) ([]entity.Pet, error) {
	clientStatus := client.FindPetsByStatusParamsStatus(status)
	params := &client.FindPetsByStatusParams{
		Status: &clientStatus,
	}

	resp, err := s.petStoreClient.FindPetsByStatusWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("network error calling petstore api: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from petstore client: %d", resp.StatusCode())
	}

	if resp.JSON200 == nil {
		return []entity.Pet{}, nil
	}

	clientPets := *resp.JSON200
	pets := make([]entity.Pet, len(clientPets))
	for i, p := range clientPets {
		pets[i] = mapClientPetToEntity(&p)
	}

	return pets, nil
}

func mapClientPetToEntity(clientPet *client.Pet) entity.Pet {
	pet := entity.Pet{
		Name:      clientPet.Name,
		PhotoURLs: clientPet.PhotoUrls,
	}

	if clientPet.Id != nil {
		pet.ID = *clientPet.Id
	}

	if clientPet.Status != nil {
		pet.Status = entity.PetStatus(*clientPet.Status)
	}

	if clientPet.Category != nil {
		if clientPet.Category.Id != nil {
			pet.Category.ID = *clientPet.Category.Id
		}
		if clientPet.Category.Name != nil {
			pet.Category.Name = *clientPet.Category.Name
		}
	}

	if clientPet.Tags != nil {
		tags := make([]string, 0, len(*clientPet.Tags))
		for _, t := range *clientPet.Tags {
			if t.Name != nil {
				tags = append(tags, *t.Name)
			}
		}
		pet.Tags = tags
	}

	return pet
}
