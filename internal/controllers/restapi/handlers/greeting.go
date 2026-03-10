package handlers

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"
)

type GreetingHandler interface {
	GetGreeting(ctx context.Context, input *models.GreetingInput) (*models.GreetingOutput, error)
}

type greetingHandler struct {
	GreetingUseCase usecases.GreetingUseCase
}

func NewGreetingHandler(greetingUseCase usecases.GreetingUseCase) *greetingHandler {
	return &greetingHandler{
		GreetingUseCase: greetingUseCase,
	}
}

func (h *greetingHandler) GetGreeting(ctx context.Context, input *models.GreetingInput) (*models.GreetingOutput, error) {
	message, err := h.GreetingUseCase.GetGreetingMessage(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	resp := models.GreetingOutput{
		Body: models.GreetingOutputBody{
			Message: message,
		},
	}

	return &resp, nil
}
