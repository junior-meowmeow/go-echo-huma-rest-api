package handlers

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"
)

type HealthHandler interface {
	GetHealthStatus(ctx context.Context, _ *struct{}) (*models.HealthOutput, error)
}

type healthHandler struct {
	HealthUseCase usecases.HealthUseCase
}

func NewHealthHandler(healthUseCase usecases.HealthUseCase) *healthHandler {
	return &healthHandler{
		HealthUseCase: healthUseCase,
	}
}

func (h *healthHandler) GetHealthStatus(ctx context.Context, _ *struct{}) (*models.HealthOutput, error) {
	status, err := h.HealthUseCase.GetHealthStatus(ctx)
	if err != nil {
		return nil, err
	}

	resp := models.HealthOutput{
		Body: models.HealthBody{
			Status: status,
		},
	}

	return &resp, nil
}
