package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

func TestResolveError(t *testing.T) {
	ctx := context.Background()

	t.Run("Should return nil when error is nil", func(t *testing.T) {
		err := handler.ResolveError(ctx, nil)
		assert.NoError(t, err)
	})

	t.Run("Should map ErrNotFound to 404", func(t *testing.T) {
		err := handler.ResolveError(ctx, entity.ErrNotFound)

		require.Error(t, err)
		require.ErrorContains(t, err, "Resource not found")
	})

	t.Run("Should map ErrAlreadyExists to 409", func(t *testing.T) {
		err := handler.ResolveError(ctx, entity.ErrAlreadyExists)

		require.Error(t, err)
		require.ErrorContains(t, err, "Resource already exists")
	})

	t.Run("Should map ErrInvalidCredentials to 401", func(t *testing.T) {
		err := handler.ResolveError(ctx, entity.ErrInvalidCredentials)

		require.Error(t, err)
		require.ErrorContains(t, err, "Invalid credentials")
	})

	t.Run("Should map unknown error to 500", func(t *testing.T) {
		err := handler.ResolveError(ctx, errors.New("random error"))

		require.Error(t, err)
		require.ErrorContains(t, err, "An unexpected internal error occurred")
	})
}
