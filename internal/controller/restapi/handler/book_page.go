package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type BookPageHandler interface {
	CreateBookPage(ctx context.Context, request *schema.CreateBookPageRequest) (*schema.CreateBookPageResponse, error)
	GetBookPages(ctx context.Context, request *schema.GetBookPagesRequest) (*schema.GetBookPagesResponse, error)
	GetBookPagesByRange(ctx context.Context, request *schema.GetBookPagesRangeRequest) (*schema.GetBookPagesResponse, error)
	GetBookPagesByOffset(ctx context.Context, request *schema.GetBookPagesOffsetRequest) (*schema.GetBookPagesResponse, error)
	GetBookPageByID(ctx context.Context, request *schema.GetBookPageByIDRequest) (*schema.GetBookPageByIDResponse, error)
}

type bookPageHandler struct {
	BookPageUseCase usecase.BookPageUseCase
}

func NewBookPageHandler(bookPageUseCase usecase.BookPageUseCase) *bookPageHandler {
	return &bookPageHandler{
		BookPageUseCase: bookPageUseCase,
	}
}

func (h *bookPageHandler) CreateBookPage(ctx context.Context, request *schema.CreateBookPageRequest) (*schema.CreateBookPageResponse, error) {
	metadata := request.Body.Metadata
	bookPage := &entity.BookPage{
		BookID:     request.Body.BookID,
		PageNumber: request.Body.PageNumber,
		Content:    request.Body.Content,
		Metadata: entity.BookPageMetadata{
			IsBookmarked: metadata.IsBookmarked,
			Highlight:    metadata.Highlight,
		},
		AttachedImageFileID: request.Body.AttachedImageFileID,
	}

	id, err := h.BookPageUseCase.CreateBookPage(ctx, bookPage)
	if err != nil {
		return nil, err
	}

	resp := schema.CreateBookPageResponse{}
	resp.Body.ID = id

	return &resp, nil
}

func (h *bookPageHandler) GetBookPages(ctx context.Context, request *schema.GetBookPagesRequest) (*schema.GetBookPagesResponse, error) {
	var bookPages []entity.BookPage
	var err error

	if request.GetAll {
		bookPages, err = h.BookPageUseCase.GetAllBookPages(ctx, request.BookID)
	} else {
		bookPages, err = h.BookPageUseCase.GetBookPagesWithPagination(ctx, request.BookID, request.PageSize, request.PageNumber)
	}

	if err != nil {
		return nil, err
	}

	bookPageOutputs := convertBookPages(bookPages)

	resp := schema.GetBookPagesResponse{}
	resp.Body.Data = bookPageOutputs

	return &resp, nil
}

func (h *bookPageHandler) GetBookPagesByRange(ctx context.Context, request *schema.GetBookPagesRangeRequest) (*schema.GetBookPagesResponse, error) {
	bookPages, err := h.BookPageUseCase.GetBookPagesByRange(ctx, request.BookID, request.StartPage, request.EndPage)
	if err != nil {
		return nil, err
	}

	bookPageOutputs := convertBookPages(bookPages)

	resp := schema.GetBookPagesResponse{}
	resp.Body.Data = bookPageOutputs

	return &resp, nil
}

func (h *bookPageHandler) GetBookPagesByOffset(ctx context.Context, request *schema.GetBookPagesOffsetRequest) (*schema.GetBookPagesResponse, error) {
	bookPages, err := h.BookPageUseCase.GetBookPagesByOffset(ctx, request.BookID, request.CenterPage, request.Offset)
	if err != nil {
		return nil, err
	}

	bookPageOutputs := convertBookPages(bookPages)

	resp := schema.GetBookPagesResponse{}
	resp.Body.Data = bookPageOutputs

	return &resp, nil
}

func (h *bookPageHandler) GetBookPageByID(ctx context.Context, request *schema.GetBookPageByIDRequest) (*schema.GetBookPageByIDResponse, error) {
	bookPage, err := h.BookPageUseCase.GetBookPageByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	bookPageOutput := convertBookPage(bookPage)

	resp := schema.GetBookPageByIDResponse{}
	resp.Body = bookPageOutput

	return &resp, nil
}

func convertBookPages(bookPages []entity.BookPage) []schema.BookPage {
	bookPageOutputs := make([]schema.BookPage, len(bookPages))
	for i, r := range bookPages {
		bookPageOutputs[i] = convertBookPage(r)
	}
	return bookPageOutputs
}

func convertBookPage(bookPage entity.BookPage) schema.BookPage {
	return schema.BookPage{
		ID:         bookPage.ID,
		BookID:     bookPage.BookID,
		PageNumber: bookPage.PageNumber,
		Content:    bookPage.Content,
		Metadata: schema.BookPageMetadata{
			IsBookmarked: bookPage.Metadata.IsBookmarked,
			Highlight:    bookPage.Metadata.Highlight,
		},
		AttachedImageFileID: bookPage.AttachedImageFileID,
	}
}
