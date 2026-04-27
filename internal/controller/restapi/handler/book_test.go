package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	usecasemocks "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase/mocks"
)

func TestBookHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateBook", func(t *testing.T) {
		t.Run("Should create book successfully", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookUseCase(t)

			mockUC.EXPECT().
				CreateBook(mock.Anything, mock.MatchedBy(func(b *entity.Book) bool {
					return b.Name == "Test Book" &&
						b.Description == "Desc" &&
						b.Metadata.Author == "Author"
				})).
				Return("book-id-123", nil)

			h := handler.NewBookHandler(mockUC)

			req := &schema.CreateBookRequest{}
			req.Body.Name = "Test Book"
			req.Body.Description = "Desc"
			req.Body.Metadata.Author = "Author"
			req.Body.Metadata.ISBN = "123"
			req.Body.Metadata.Genre = "Fiction"
			req.Body.CoverImageFileID = "file-id"

			resp, err := h.CreateBook(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, "book-id-123", resp.Body.ID)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookUseCase(t)

			mockUC.EXPECT().
				CreateBook(mock.Anything, mock.Anything).
				Return("", errors.New("db error"))

			h := handler.NewBookHandler(mockUC)

			req := &schema.CreateBookRequest{}

			resp, err := h.CreateBook(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetBooks", func(t *testing.T) {
		mockBooks := []entity.Book{
			{
				ID:          "1",
				Name:        "Book 1",
				Description: "Desc 1",
				Metadata: entity.BookMetadata{
					Author: "Author 1",
					ISBN:   "111",
					Genre:  "Fiction",
				},
				CoverImageFileID: "file1",
				CreatedAt:        time.Now(),
			},
			{
				ID:          "2",
				Name:        "Book 2",
				Description: "Desc 2",
				Metadata: entity.BookMetadata{
					Author: "Author 2",
					ISBN:   "222",
					Genre:  "Drama",
				},
				CoverImageFileID: "file2",
				CreatedAt:        time.Now(),
			},
		}

		t.Run("Should return all books when GetAll is true", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookUseCase(t)

			mockUC.EXPECT().
				GetAllBooks(mock.Anything).
				Return(mockBooks, nil)

			h := handler.NewBookHandler(mockUC)

			req := &schema.GetBooksRequest{
				GetAll: true,
			}

			resp, err := h.GetBooks(ctx, req)

			require.NoError(t, err)
			require.Len(t, resp.Body.Data, 2)

			assert.Equal(t, "Book 1", resp.Body.Data[0].Name)
			assert.Equal(t, "Author 1", resp.Body.Data[0].Metadata.Author)
		})

		t.Run("Should return paginated books when GetAll is false", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookUseCase(t)

			mockUC.EXPECT().
				GetBooksWithPagination(mock.Anything, int64(2), int64(1)).
				Return(mockBooks, nil)

			h := handler.NewBookHandler(mockUC)

			req := &schema.GetBooksRequest{
				GetAll:     false,
				PageSize:   2,
				PageNumber: 1,
			}

			resp, err := h.GetBooks(ctx, req)

			require.NoError(t, err)
			require.Len(t, resp.Body.Data, 2)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookUseCase(t)

			mockUC.EXPECT().
				GetAllBooks(mock.Anything).
				Return(nil, errors.New("error"))

			h := handler.NewBookHandler(mockUC)

			req := &schema.GetBooksRequest{
				GetAll: true,
			}

			resp, err := h.GetBooks(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetBookByID", func(t *testing.T) {
		t.Run("Should return mapped book", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookUseCase(t)

			now := time.Now()

			book := entity.Book{
				ID:          "1",
				Name:        "Book 1",
				Description: "Desc",
				Metadata: entity.BookMetadata{
					Author: "Author",
					ISBN:   "123",
					Genre:  "Fiction",
				},
				CoverImageFileID: "file-id",
				CreatedAt:        now,
			}

			mockUC.EXPECT().
				GetBookByID(mock.Anything, "1").
				Return(book, nil)

			h := handler.NewBookHandler(mockUC)

			req := &schema.GetBookByIDRequest{
				ID: "1",
			}

			resp, err := h.GetBookByID(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, "Book 1", resp.Body.Name)
			assert.Equal(t, "Author", resp.Body.Metadata.Author)
			assert.Equal(t, "file-id", resp.Body.CoverImageFileID)
			assert.Equal(t, now, resp.Body.CreatedAt)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookUseCase(t)

			mockUC.EXPECT().
				GetBookByID(mock.Anything, "1").
				Return(entity.Book{}, errors.New("not found"))

			h := handler.NewBookHandler(mockUC)

			req := &schema.GetBookByIDRequest{
				ID: "1",
			}

			resp, err := h.GetBookByID(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})
}
