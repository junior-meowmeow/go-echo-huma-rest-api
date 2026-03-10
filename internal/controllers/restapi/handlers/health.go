package handlers

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
)

type HealthHandler interface {
	GetHealthStatus(ctx context.Context, _ *struct{}) (*models.HealthOutput, error)
}

type healthHandler struct {
}

func NewHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) GetHealthStatus(ctx context.Context, _ *struct{}) (*models.HealthOutput, error) {
	resp := &models.HealthOutput{
		Body: models.HealthBody{
			Status: "ok",
		},
	}

	return resp, nil
}
