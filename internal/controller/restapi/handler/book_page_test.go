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

func TestBookPageHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("CreateBookPage", func(t *testing.T) {
		t.Run("Should create book page successfully", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				CreateBookPage(mock.Anything, mock.MatchedBy(func(p *entity.BookPage) bool {
					return p.BookID == "book1" &&
						p.PageNumber == 1 &&
						p.Content == "content" &&
						p.Metadata.IsBookmarked &&
						p.Metadata.Highlight == "highlight"
				})).
				Return("page-id-123", nil)

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.CreateBookPageRequest{}
			req.Body.BookID = "book1"
			req.Body.PageNumber = 1
			req.Body.Content = "content"
			req.Body.Metadata.IsBookmarked = true
			req.Body.Metadata.Highlight = "highlight"
			req.Body.AttachedImageFileID = "file1"

			resp, err := h.CreateBookPage(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, "page-id-123", resp.Body.ID)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				CreateBookPage(mock.Anything, mock.Anything).
				Return("", errors.New("fail"))

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.CreateBookPageRequest{}

			resp, err := h.CreateBookPage(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetBookPages", func(t *testing.T) {
		mockPages := []entity.BookPage{
			{ID: "1", BookID: "book1", PageNumber: 1, Content: "p1"},
			{ID: "2", BookID: "book1", PageNumber: 2, Content: "p2"},
		}

		t.Run("Should return all pages when GetAll is true", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetAllBookPages(mock.Anything, "book1").
				Return(mockPages, nil)

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPagesRequest{
				ParentBookIDQuery: schema.ParentBookIDQuery{
					BookID: "book1",
				},
				GetAll: true,
			}

			resp, err := h.GetBookPages(ctx, req)

			require.NoError(t, err)
			require.Len(t, resp.Body.Data, 2)
			assert.Equal(t, "1", resp.Body.Data[0].ID)
		})

		t.Run("Should return paginated pages when GetAll is false", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetBookPagesWithPagination(mock.Anything, "book1", int64(2), int64(1)).
				Return(mockPages, nil)

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPagesRequest{
				ParentBookIDQuery: schema.ParentBookIDQuery{
					BookID: "book1",
				},
				GetAll:     false,
				PageSize:   2,
				PageNumber: 1,
			}

			resp, err := h.GetBookPages(ctx, req)

			require.NoError(t, err)
			require.Len(t, resp.Body.Data, 2)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetAllBookPages(mock.Anything, "book1").
				Return(nil, errors.New("fail"))

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPagesRequest{
				ParentBookIDQuery: schema.ParentBookIDQuery{
					BookID: "book1",
				},
				GetAll: true,
			}

			resp, err := h.GetBookPages(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetBookPagesByRange", func(t *testing.T) {
		mockPages := []entity.BookPage{
			{ID: "1", PageNumber: 1},
			{ID: "2", PageNumber: 2},
		}

		t.Run("Should return pages in range", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetBookPagesByRange(mock.Anything, "book1", int64(1), int64(2)).
				Return(mockPages, nil)

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPagesRangeRequest{
				ParentBookIDQuery: schema.ParentBookIDQuery{
					BookID: "book1",
				},
				StartPage: 1,
				EndPage:   2,
			}

			resp, err := h.GetBookPagesByRange(ctx, req)

			require.NoError(t, err)
			require.Len(t, resp.Body.Data, 2)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetBookPagesByRange(mock.Anything, "book1", int64(1), int64(2)).
				Return(nil, errors.New("fail"))

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPagesRangeRequest{
				ParentBookIDQuery: schema.ParentBookIDQuery{
					BookID: "book1",
				},
				StartPage: 1,
				EndPage:   2,
			}

			resp, err := h.GetBookPagesByRange(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetBookPagesByOffset", func(t *testing.T) {
		mockPages := []entity.BookPage{
			{ID: "1"}, {ID: "2"}, {ID: "3"},
		}

		t.Run("Should return pages by offset", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetBookPagesByOffset(mock.Anything, "book1", int64(5), int64(2)).
				Return(mockPages, nil)

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPagesOffsetRequest{
				ParentBookIDQuery: schema.ParentBookIDQuery{
					BookID: "book1",
				},
				CenterPage: 5,
				Offset:     2,
			}

			resp, err := h.GetBookPagesByOffset(ctx, req)

			require.NoError(t, err)
			require.Len(t, resp.Body.Data, 3)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetBookPagesByOffset(mock.Anything, "book1", int64(5), int64(2)).
				Return(nil, errors.New("fail"))

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPagesOffsetRequest{
				ParentBookIDQuery: schema.ParentBookIDQuery{
					BookID: "book1",
				},
				CenterPage: 5,
				Offset:     2,
			}

			resp, err := h.GetBookPagesByOffset(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})

	t.Run("GetBookPageByID", func(t *testing.T) {
		t.Run("Should return mapped page", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			page := entity.BookPage{
				ID:         "1",
				BookID:     "book1",
				PageNumber: 10,
				Content:    "hello",
			}

			mockUC.EXPECT().
				GetBookPageByID(mock.Anything, "1").
				Return(page, nil)

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPageByIDRequest{
				ID: "1",
			}

			resp, err := h.GetBookPageByID(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, "1", resp.Body.ID)
			assert.Equal(t, "book1", resp.Body.BookID)
			assert.Equal(t, int64(10), resp.Body.PageNumber)
		})

		t.Run("Should return error when usecase fails", func(t *testing.T) {
			mockUC := usecasemocks.NewMockBookPageUseCase(t)

			mockUC.EXPECT().
				GetBookPageByID(mock.Anything, "1").
				Return(entity.BookPage{}, errors.New("fail"))

			h := handler.NewBookPageHandler(mockUC)

			req := &schema.GetBookPageByIDRequest{
				ID: "1",
			}

			resp, err := h.GetBookPageByID(ctx, req)

			require.Error(t, err)
			assert.Nil(t, resp)
			require.ErrorContains(t, err, "An unexpected internal error occurred")
		})
	})
}
