package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

func TestGreetingUseCase(t *testing.T) {
	ctx := context.Background()
	uc := usecase.NewGreetingUseCase()

	t.Run("Should return greeting message", func(t *testing.T) {
		result := uc.GetGreetingMessage(ctx, "John")
		assert.Equal(t, "Hello, John!", result)
	})

	t.Run("Should handle empty name", func(t *testing.T) {
		result := uc.GetGreetingMessage(ctx, "")
		assert.Equal(t, "Hello, !", result)
	})
}
