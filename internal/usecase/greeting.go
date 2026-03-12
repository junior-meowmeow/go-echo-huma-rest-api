package usecase

import (
	"context"
	"fmt"
)

type GreetingUseCase interface {
	GetGreetingMessage(ctx context.Context, name string) (string, error)
}

type greetingUseCase struct {
}

func NewGreetingUseCase() *greetingUseCase {
	return &greetingUseCase{}
}

func (u *greetingUseCase) GetGreetingMessage(ctx context.Context, name string) (string, error) {
	message := fmt.Sprintf("Hello, %s!", name)

	return message, nil
}
