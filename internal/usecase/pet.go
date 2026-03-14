package usecase

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/external"
)

type PetUseCase interface {
	GetAvailablePets(ctx context.Context) ([]entity.Pet, error)
	GetPetByID(ctx context.Context, id int64) (entity.Pet, error)
}

type petUseCase struct {
	PetService external.PetService
}

func NewPetUseCase(petService external.PetService) *petUseCase {
	return &petUseCase{
		PetService: petService,
	}
}

func (u *petUseCase) GetAvailablePets(ctx context.Context) ([]entity.Pet, error) {
	pets, err := u.PetService.GetPetsByStatus(ctx, entity.PetStatusAvailable)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available pets: %w", err)
	}

	return pets, nil
}

func (u *petUseCase) GetPetByID(ctx context.Context, id int64) (entity.Pet, error) {
	pet, err := u.PetService.GetPetByID(ctx, id)
	if err != nil {
		return entity.Pet{}, fmt.Errorf("failed to fetch pet: %w", err)
	}

	return pet, nil
}
