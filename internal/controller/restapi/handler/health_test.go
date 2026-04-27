package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	usecasemocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase/mocks"
)

func TestHealthHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetHealthStatus", func(t *testing.T) {
		t.Run("Should return health status successfully", func(t *testing.T) {
			mockUC := usecasemocks.NewMockHealthUseCase(t)

			mockUC.EXPECT().
				GetHealthStatus(mock.Anything).
				Return("ok", nil)

			h := handler.NewHealthHandler(mockUC)

			resp, err := h.GetHealthStatus(ctx, &schema.GetHealthStatusRequest{})

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, "ok", resp.Body.Status)
		})

		t.Run("Should map error to 500 when unknown error", func(t *testing.T) {
			mockUC := usecasemocks.NewMockHealthUseCase(t)

			mockUC.EXPECT().
				GetHealthStatus(mock.Anything).
				Return("", errors.New("some internal error"))

			h := handler.NewHealthHandler(mockUC)

			resp, err := h.GetHealthStatus(ctx, &schema.GetHealthStatusRequest{})

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})
}
