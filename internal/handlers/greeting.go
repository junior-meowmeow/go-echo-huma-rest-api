package handlers

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/models"
)

func (h *Handler) GetGreeting(ctx context.Context, input *models.GreetingInput) (*models.GreetingOutput, error) {
	resp := &models.GreetingOutput{
		Body: models.GreetingOutputBody{
			Message: fmt.Sprintf("Hello, %s!", input.Name),
		},
	}

	return resp, nil
}
