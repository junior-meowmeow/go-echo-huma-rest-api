package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port/mocks"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

func TestBookPageUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateBookPage", func(t *testing.T) {
		t.Run("Should create book page successfully", func(t *testing.T) {
			mockBookRepo := mocks.NewMockBookRepository(t)
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)

			input := &entity.BookPage{
				BookID:     "book-1",
				PageNumber: 1,
			}

			mockBookRepo.EXPECT().
				GetBookByID(ctx, "book-1").
				Return(entity.Book{ID: "book-1"}, nil)

			mockBookPageRepo.EXPECT().
				CreateBookPage(ctx, input).
				Run(func(_ context.Context, bp *entity.BookPage) {
					assert.False(t, bp.CreatedAt.IsZero())
					assert.False(t, bp.UpdatedAt.IsZero())
					assert.WithinDuration(t, bp.CreatedAt, bp.UpdatedAt, time.Millisecond)
				}).
				Return("page-id", nil)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			id, err := uc.CreateBookPage(ctx, input)

			require.NoError(t, err)
			assert.Equal(t, "page-id", id)
		})

		t.Run("Should return error when book not found", func(t *testing.T) {
			mockBookRepo := mocks.NewMockBookRepository(t)
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)

			mockErr := errors.New("not found")

			mockBookRepo.EXPECT().
				GetBookByID(ctx, "book-1").
				Return(entity.Book{}, mockErr)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			_, err := uc.CreateBookPage(ctx, &entity.BookPage{BookID: "book-1"})

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to fetch book")
			assert.ErrorIs(t, err, mockErr)
		})

		t.Run("Should return error when repository fails", func(t *testing.T) {
			mockBookRepo := mocks.NewMockBookRepository(t)
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)

			mockErr := errors.New("insert error")

			mockBookRepo.EXPECT().
				GetBookByID(ctx, "book-1").
				Return(entity.Book{}, nil)

			mockBookPageRepo.EXPECT().
				CreateBookPage(ctx, mock.Anything).
				Return("", mockErr)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			_, err := uc.CreateBookPage(ctx, &entity.BookPage{BookID: "book-1"})

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to create book page")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetAllBookPages", func(t *testing.T) {
		t.Run("Should return book pages", func(t *testing.T) {
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)
			mockBookRepo := mocks.NewMockBookRepository(t)

			expected := []entity.BookPage{{PageNumber: 1}}

			mockBookPageRepo.EXPECT().
				GetBookPagesByBookID(ctx, "book-1").
				Return(expected, nil)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			result, err := uc.GetAllBookPages(ctx, "book-1")

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("GetBookPagesWithPagination", func(t *testing.T) {
		t.Run("Should return book pages", func(t *testing.T) {
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)
			mockBookRepo := mocks.NewMockBookRepository(t)

			expected := []entity.BookPage{{PageNumber: 2}}

			mockBookPageRepo.EXPECT().
				GetBookpagesByBookIDWithPagination(ctx, "book-1", int64(2), int64(1)).
				Return(expected, nil)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			result, err := uc.GetBookPagesWithPagination(ctx, "book-1", 2, 1)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("GetBookPagesByRange", func(t *testing.T) {
		t.Run("Should return book pages", func(t *testing.T) {
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)
			mockBookRepo := mocks.NewMockBookRepository(t)

			expected := []entity.BookPage{{PageNumber: 3}}

			mockBookPageRepo.EXPECT().
				GetBookpagesByPageRange(ctx, "book-1", int64(2), int64(4)).
				Return(expected, nil)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			result, err := uc.GetBookPagesByRange(ctx, "book-1", 2, 4)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("GetBookPagesByOffset", func(t *testing.T) {
		t.Run("Should return book pages", func(t *testing.T) {
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)
			mockBookRepo := mocks.NewMockBookRepository(t)

			expected := []entity.BookPage{{PageNumber: 5}}

			mockBookPageRepo.EXPECT().
				GetBookpagesAroundPageNumber(ctx, "book-1", int64(5), int64(1)).
				Return(expected, nil)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			result, err := uc.GetBookPagesByOffset(ctx, "book-1", 5, 1)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	})

	t.Run("GetBookPageByID", func(t *testing.T) {
		t.Run("Should return book page", func(t *testing.T) {
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)
			mockBookRepo := mocks.NewMockBookRepository(t)

			expected := entity.BookPage{PageNumber: 1}

			mockBookPageRepo.EXPECT().
				GetBookPageByID(ctx, "page-1").
				Return(expected, nil)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			result, err := uc.GetBookPageByID(ctx, "page-1")

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return wrapped error", func(t *testing.T) {
			mockBookPageRepo := mocks.NewMockBookPageRepository(t)
			mockBookRepo := mocks.NewMockBookRepository(t)

			mockErr := errors.New("db error")

			mockBookPageRepo.EXPECT().
				GetBookPageByID(ctx, "page-1").
				Return(entity.BookPage{}, mockErr)

			uc := usecase.NewBookPageUseCase(mockBookRepo, mockBookPageRepo)

			_, err := uc.GetBookPageByID(ctx, "page-1")

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to fetch book page")
			assert.ErrorIs(t, err, mockErr)
		})
	})
}
