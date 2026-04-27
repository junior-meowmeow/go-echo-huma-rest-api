package handler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	usecasemocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase/mocks"
)

func TestGreetingHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetGreeting", func(t *testing.T) {
		t.Run("Should return greeting response", func(t *testing.T) {
			mockUC := usecasemocks.NewMockGreetingUseCase(t)

			req := &schema.GreetingRequest{
				Name: "John",
			}

			mockUC.EXPECT().
				GetGreetingMessage(ctx, "John").
				Return("Hello, John!")

			h := handler.NewGreetingHandler(mockUC)

			resp, err := h.GetGreeting(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, "Hello, John!", resp.Body.Message)
		})

		t.Run("Should handle empty name", func(t *testing.T) {
			mockUC := usecasemocks.NewMockGreetingUseCase(t)

			req := &schema.GreetingRequest{
				Name: "",
			}

			mockUC.EXPECT().
				GetGreetingMessage(ctx, "").
				Return("Hello, !")

			h := handler.NewGreetingHandler(mockUC)

			resp, err := h.GetGreeting(ctx, req)

			require.NoError(t, err)
			assert.Equal(t, "Hello, !", resp.Body.Message)
		})
	})
}
