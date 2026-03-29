package usecase

import (
	"context"
)

type HealthUseCase interface {
	GetHealthStatus(ctx context.Context) (string, error)
}

type healthUseCase struct {
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewHealthUseCase() *healthUseCase {
	return &healthUseCase{}
}

//revive:enable:unexported-return

func (u *healthUseCase) GetHealthStatus(_ context.Context) (string, error) {
	// Only check server for now
	// It is possible to add external service checks later e.g., ping database, etc.

	return "ok", nil
}
