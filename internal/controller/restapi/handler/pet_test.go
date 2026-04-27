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
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	usecasemocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase/mocks"
)

func TestPetHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("GetAvailablePets", func(t *testing.T) {
		t.Run("Should return mapped pets", func(t *testing.T) {
			mockUC := usecasemocks.NewMockPetUseCase(t)

			input := []entity.Pet{
				{
					ID:   1,
					Name: "Dog",
					Category: entity.PetCategory{
						ID:   10,
						Name: "Doggy",
					},
					PhotoURLs: []string{"url1"},
					Status:    entity.PetStatusAvailable,
					Tags:      []string{"friendly"},
				},
				{
					ID:   2,
					Name: "Cat",
					Category: entity.PetCategory{
						ID:   20,
						Name: "Catty",
					},
					PhotoURLs: []string{"url2"},
					Status:    entity.PetStatusAvailable,
					Tags:      []string{"cute"},
				},
			}

			mockUC.EXPECT().
				GetAvailablePets(mock.Anything).
				Return(input, nil)

			h := handler.NewPetHandler(mockUC)

			resp, err := h.GetAvailablePets(ctx, &schema.GetAvailablePetsRequest{})

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp.Body.Data, 2)

			assert.Equal(t, int64(1), resp.Body.Data[0].ID)
			assert.Equal(t, "Dog", resp.Body.Data[0].Name)
			assert.Equal(t, int64(10), resp.Body.Data[0].Category.ID)
			assert.Equal(t, "Doggy", resp.Body.Data[0].Category.Name)
			assert.Equal(t, []string{"url1"}, resp.Body.Data[0].PhotoURLs)
			assert.Equal(t, "available", resp.Body.Data[0].Status)
			assert.Equal(t, []string{"friendly"}, resp.Body.Data[0].Tags)

			assert.Equal(t, int64(2), resp.Body.Data[1].ID)
			assert.Equal(t, "Cat", resp.Body.Data[1].Name)
		})

		t.Run("Should return empty list", func(t *testing.T) {
			mockUC := usecasemocks.NewMockPetUseCase(t)

			mockUC.EXPECT().
				GetAvailablePets(mock.Anything).
				Return([]entity.Pet{}, nil)

			h := handler.NewPetHandler(mockUC)

			resp, err := h.GetAvailablePets(ctx, &schema.GetAvailablePetsRequest{})

			require.NoError(t, err)
			assert.Empty(t, resp.Body.Data)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockPetUseCase(t)

			mockErr := errors.New("service error")

			mockUC.EXPECT().
				GetAvailablePets(mock.Anything).
				Return(nil, mockErr)

			h := handler.NewPetHandler(mockUC)

			resp, err := h.GetAvailablePets(ctx, &schema.GetAvailablePetsRequest{})

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetPetByID", func(t *testing.T) {
		t.Run("Should return mapped pet", func(t *testing.T) {
			mockUC := usecasemocks.NewMockPetUseCase(t)

			input := entity.Pet{
				ID:   1,
				Name: "Dog",
				Category: entity.PetCategory{
					ID:   10,
					Name: "Doggy",
				},
				PhotoURLs: []string{"url1"},
				Status:    entity.PetStatusAvailable,
				Tags:      []string{"friendly"},
			}

			mockUC.EXPECT().
				GetPetByID(mock.Anything, int64(1)).
				Return(input, nil)

			h := handler.NewPetHandler(mockUC)

			req := &schema.GetPetByIDRequest{ID: 1}

			resp, err := h.GetPetByID(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, int64(1), resp.Body.ID)
			assert.Equal(t, "Dog", resp.Body.Name)
			assert.Equal(t, int64(10), resp.Body.Category.ID)
			assert.Equal(t, "Doggy", resp.Body.Category.Name)
			assert.Equal(t, []string{"url1"}, resp.Body.PhotoURLs)
			assert.Equal(t, "available", resp.Body.Status)
			assert.Equal(t, []string{"friendly"}, resp.Body.Tags)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockPetUseCase(t)

			mockErr := errors.New("not found")

			mockUC.EXPECT().
				GetPetByID(mock.Anything, int64(1)).
				Return(entity.Pet{}, mockErr)

			h := handler.NewPetHandler(mockUC)

			req := &schema.GetPetByIDRequest{ID: 1}

			resp, err := h.GetPetByID(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)

			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})
}
