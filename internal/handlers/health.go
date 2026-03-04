package handlers

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/models"
)

func (h *Handler) GetHealthStatus(ctx context.Context, _ *struct{}) (*models.HealthOutput, error) {
	resp := &models.HealthOutput{
		Body: models.HealthBody{
			Status: "ok",
		},
	}

	return resp, nil
}
