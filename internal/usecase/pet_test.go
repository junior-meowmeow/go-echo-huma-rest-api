package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port/mocks"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

func TestPetUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("GetAvailablePets", func(t *testing.T) {
		t.Run("Should return available pets", func(t *testing.T) {
			mockService := mocks.NewMockPetService(t)

			expected := []entity.Pet{
				{ID: 1, Name: "Dog"},
				{ID: 2, Name: "Cat"},
			}

			mockService.EXPECT().
				GetPetsByStatus(ctx, entity.PetStatusAvailable).
				Return(expected, nil)

			uc := usecase.NewPetUseCase(mockService)

			result, err := uc.GetAvailablePets(ctx)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return wrapped error when service fails", func(t *testing.T) {
			mockService := mocks.NewMockPetService(t)

			mockErr := errors.New("service error")

			mockService.EXPECT().
				GetPetsByStatus(ctx, entity.PetStatusAvailable).
				Return(nil, mockErr)

			uc := usecase.NewPetUseCase(mockService)

			result, err := uc.GetAvailablePets(ctx)

			require.Error(t, err)
			assert.Nil(t, result)
			require.ErrorContains(t, err, "failed to fetch available pets")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetPetByID", func(t *testing.T) {
		t.Run("Should return pet when exists", func(t *testing.T) {
			mockService := mocks.NewMockPetService(t)

			expected := entity.Pet{
				ID:   1,
				Name: "Dog",
			}

			mockService.EXPECT().
				GetPetByID(ctx, int64(1)).
				Return(expected, nil)

			uc := usecase.NewPetUseCase(mockService)

			result, err := uc.GetPetByID(ctx, 1)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return wrapped error when service fails", func(t *testing.T) {
			mockService := mocks.NewMockPetService(t)

			mockErr := errors.New("service error")

			mockService.EXPECT().
				GetPetByID(ctx, int64(1)).
				Return(entity.Pet{}, mockErr)

			uc := usecase.NewPetUseCase(mockService)

			result, err := uc.GetPetByID(ctx, 1)

			require.Error(t, err)
			assert.Equal(t, entity.Pet{}, result)
			require.ErrorContains(t, err, "failed to fetch pet")
			assert.ErrorIs(t, err, mockErr)
		})
	})
}
