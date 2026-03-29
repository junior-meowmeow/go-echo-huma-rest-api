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

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewGreetingUseCase() *greetingUseCase {
	return &greetingUseCase{}
}

//revive:enable:unexported-return

func (u *greetingUseCase) GetGreetingMessage(_ context.Context, name string) string {
	message := fmt.Sprintf("Hello, %s!", name)

	return message
}
