package usecases

import (
	"context"
)

type HealthUseCase interface {
	GetHealthStatus(ctx context.Context) (string, error)
}

type healthUseCase struct {
}

func NewHealthUseCase() *healthUseCase {
	return &healthUseCase{}
}

func (u *healthUseCase) GetHealthStatus(ctx context.Context) (string, error) {
	// Only check server for now
	// It is possible to add external service checks later e.g., ping database, etc.

	return "ok", nil
}
