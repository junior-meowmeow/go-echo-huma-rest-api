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

func TestBookUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateBook", func(t *testing.T) {
		t.Run("Should create book successfully", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			input := &entity.Book{
				Name: "Test Book",
			}

			mockRepo.EXPECT().
				CreateBook(ctx, input).
				Run(func(_ context.Context, b *entity.Book) {
					assert.False(t, b.CreatedAt.IsZero())
					assert.False(t, b.UpdatedAt.IsZero())
					assert.WithinDuration(t, b.CreatedAt, b.UpdatedAt, time.Millisecond)
				}).
				Return("book-id", nil)

			uc := usecase.NewBookUseCase(mockRepo)

			id, err := uc.CreateBook(ctx, input)

			require.NoError(t, err)
			assert.Equal(t, "book-id", id)
		})

		t.Run("Should return error when repository fails", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			mockErr := errors.New("insert error")

			mockRepo.EXPECT().
				CreateBook(ctx, mock.Anything).
				Return("", mockErr)

			uc := usecase.NewBookUseCase(mockRepo)

			_, err := uc.CreateBook(ctx, &entity.Book{Name: "Test"})

			require.Error(t, err)
			require.ErrorContains(t, err, "failed to create book")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetAllBooks", func(t *testing.T) {
		t.Run("Should return books", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			expected := []entity.Book{
				{Name: "Book 1"},
				{Name: "Book 2"},
			}

			mockRepo.EXPECT().
				GetAllBooks(ctx).
				Return(expected, nil)

			uc := usecase.NewBookUseCase(mockRepo)

			result, err := uc.GetAllBooks(ctx)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return wrapped error", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			mockErr := errors.New("db error")

			mockRepo.EXPECT().
				GetAllBooks(ctx).
				Return(nil, mockErr)

			uc := usecase.NewBookUseCase(mockRepo)

			result, err := uc.GetAllBooks(ctx)

			require.Error(t, err)
			assert.Nil(t, result)
			require.ErrorContains(t, err, "failed to fetch books")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetBooksWithPagination", func(t *testing.T) {
		t.Run("Should return books", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			expected := []entity.Book{
				{Name: "Book A"},
			}

			mockRepo.EXPECT().
				GetBooksWithPagination(ctx, int64(2), int64(1)).
				Return(expected, nil)

			uc := usecase.NewBookUseCase(mockRepo)

			result, err := uc.GetBooksWithPagination(ctx, 2, 1)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return wrapped error", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			mockErr := errors.New("db error")

			mockRepo.EXPECT().
				GetBooksWithPagination(ctx, int64(2), int64(1)).
				Return(nil, mockErr)

			uc := usecase.NewBookUseCase(mockRepo)

			result, err := uc.GetBooksWithPagination(ctx, 2, 1)

			require.Error(t, err)
			assert.Nil(t, result)
			require.ErrorContains(t, err, "failed to fetch books")
			assert.ErrorIs(t, err, mockErr)
		})
	})

	t.Run("GetBookByID", func(t *testing.T) {
		t.Run("Should return book", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			expected := entity.Book{
				ID:   "1",
				Name: "Book",
			}

			mockRepo.EXPECT().
				GetBookByID(ctx, "1").
				Return(expected, nil)

			uc := usecase.NewBookUseCase(mockRepo)

			result, err := uc.GetBookByID(ctx, "1")

			require.NoError(t, err)
			assert.Equal(t, expected, result)
		})

		t.Run("Should return wrapped error", func(t *testing.T) {
			mockRepo := mocks.NewMockBookRepository(t)

			mockErr := errors.New("db error")

			mockRepo.EXPECT().
				GetBookByID(ctx, "1").
				Return(entity.Book{}, mockErr)

			uc := usecase.NewBookUseCase(mockRepo)

			result, err := uc.GetBookByID(ctx, "1")

			require.Error(t, err)
			assert.Equal(t, entity.Book{}, result)
			require.ErrorContains(t, err, "failed to fetch book")
			assert.ErrorIs(t, err, mockErr)
		})
	})
}
