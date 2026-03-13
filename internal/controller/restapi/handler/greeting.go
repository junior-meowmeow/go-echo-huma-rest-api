package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type GreetingHandler interface {
	GetGreeting(ctx context.Context, request *schema.GreetingRequest) (*schema.GreetingResponse, error)
}

type greetingHandler struct {
	GreetingUseCase usecase.GreetingUseCase
}

func NewGreetingHandler(greetingUseCase usecase.GreetingUseCase) *greetingHandler {
	return &greetingHandler{
		GreetingUseCase: greetingUseCase,
	}
}

func (h *greetingHandler) GetGreeting(ctx context.Context, request *schema.GreetingRequest) (*schema.GreetingResponse, error) {
	message, err := h.GreetingUseCase.GetGreetingMessage(ctx, request.Name)
	if err != nil {
		return nil, err
	}

	resp := schema.GreetingResponse{}
	resp.Body.Message = message

	return &resp, nil
}
