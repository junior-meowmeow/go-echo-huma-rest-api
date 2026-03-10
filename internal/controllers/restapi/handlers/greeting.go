package handlers

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
)

type GreetingHandler interface {
	GetGreeting(ctx context.Context, input *models.GreetingInput) (*models.GreetingOutput, error)
}

type greetingHandler struct {
}

func NewGreetingHandler() *greetingHandler {
	return &greetingHandler{}
}

func (h *greetingHandler) GetGreeting(ctx context.Context, input *models.GreetingInput) (*models.GreetingOutput, error) {
	resp := &models.GreetingOutput{
		Body: models.GreetingOutputBody{
			Message: fmt.Sprintf("Hello, %s!", input.Name),
		},
	}

	return resp, nil
}
