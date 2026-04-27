package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

func TestHealthUseCase(t *testing.T) {
	ctx := context.Background()
	uc := usecase.NewHealthUseCase()

	t.Run("Should return ok status", func(t *testing.T) {
		status, err := uc.GetHealthStatus(ctx)

		require.NoError(t, err)
		assert.Equal(t, "ok", status)
	})
}
