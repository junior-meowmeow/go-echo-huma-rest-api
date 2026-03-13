package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type HealthHandler interface {
	GetHealthStatus(ctx context.Context, request *schema.GetHealthStatusRequest) (*schema.GetHealthStatusResponse, error)
}

type healthHandler struct {
	HealthUseCase usecase.HealthUseCase
}

func NewHealthHandler(healthUseCase usecase.HealthUseCase) *healthHandler {
	return &healthHandler{
		HealthUseCase: healthUseCase,
	}
}

func (h *healthHandler) GetHealthStatus(ctx context.Context, _ *schema.GetHealthStatusRequest) (*schema.GetHealthStatusResponse, error) {
	status, err := h.HealthUseCase.GetHealthStatus(ctx)
	if err != nil {
		return nil, err
	}

	resp := schema.GetHealthStatusResponse{}
	resp.Body.Status = status

	return &resp, nil
}
