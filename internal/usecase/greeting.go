package usecase

import (
	"context"
	"fmt"
)

type GreetingUseCase interface {
	GetGreetingMessage(ctx context.Context, name string) string
}

type greetingUseCase struct {
}

func NewGreetingUseCase() *greetingUseCase {
	return &greetingUseCase{}
}

func (u *greetingUseCase) GetGreetingMessage(ctx context.Context, name string) string {
	message := fmt.Sprintf("Hello, %s!", name)

	return message
}
