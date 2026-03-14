package external

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
)

type PetService interface {
	GetPetByID(ctx context.Context, id int64) (entity.Pet, error)
	GetPetsByStatus(ctx context.Context, status entity.PetStatus) ([]entity.Pet, error)
}
