package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type BookPageHandler interface {
	CreateBookPage(ctx context.Context, input *schema.CreateBookPageInput) (*schema.CreateBookPageOutput, error)
	GetBookPages(ctx context.Context, input *schema.GetBookPagesInput) (*schema.GetBookPagesOutput, error)
	GetBookPagesByRange(ctx context.Context, input *schema.GetBookPagesRangeInput) (*schema.GetBookPagesOutput, error)
	GetBookPagesByOffset(ctx context.Context, input *schema.GetBookPagesOffsetInput) (*schema.GetBookPagesOutput, error)
	GetBookPageByID(ctx context.Context, input *schema.GetBookPageByIDInput) (*schema.GetBookPageByIDOutput, error)
}

type bookPageHandler struct {
	BookPageUseCase usecase.BookPageUseCase
}

func NewBookPageHandler(bookPageUseCase usecase.BookPageUseCase) *bookPageHandler {
	return &bookPageHandler{
		BookPageUseCase: bookPageUseCase,
	}
}

func (h *bookPageHandler) CreateBookPage(ctx context.Context, input *schema.CreateBookPageInput) (*schema.CreateBookPageOutput, error) {
	metadata := input.Body.Metadata
	bookPage := &entity.BookPage{
		BookID:     input.Body.BookID,
		PageNumber: input.Body.PageNumber,
		Content:    input.Body.Content,
		Metadata: entity.BookPageMetadata{
			IsBookmarked: metadata.IsBookmarked,
			Highlight:    metadata.Highlight,
		},
		AttachedImageFileID: input.Body.AttachedImageFileID,
	}

	id, err := h.BookPageUseCase.CreateBookPage(ctx, bookPage)
	if err != nil {
		return nil, err
	}

	resp := schema.CreateBookPageOutput{
		Body: schema.CreateBookPageOutputBody{
			ID: id,
		},
	}

	return &resp, nil
}

func (h *bookPageHandler) GetBookPages(ctx context.Context, input *schema.GetBookPagesInput) (*schema.GetBookPagesOutput, error) {
	var bookPages []entity.BookPage
	var err error

	if input.GetAll {
		bookPages, err = h.BookPageUseCase.GetAllBookPages(ctx, input.BookID)
	} else {
		bookPages, err = h.BookPageUseCase.GetBookPagesWithPagination(ctx, input.BookID, input.PageSize, input.PageNumber)
	}

	if err != nil {
		return nil, err
	}

	bookPagesOutput := convertBookPages(bookPages)

	resp := schema.GetBookPagesOutput{
		Body: schema.GetBookPagesOutputBody{
			Data: bookPagesOutput,
		},
	}

	return &resp, nil
}

func (h *bookPageHandler) GetBookPagesByRange(ctx context.Context, input *schema.GetBookPagesRangeInput) (*schema.GetBookPagesOutput, error) {
	bookPages, err := h.BookPageUseCase.GetBookPagesByRange(ctx, input.BookID, input.StartPage, input.EndPage)
	if err != nil {
		return nil, err
	}

	bookPageOutputs := convertBookPages(bookPages)

	resp := schema.GetBookPagesOutput{
		Body: schema.GetBookPagesOutputBody{
			Data: bookPageOutputs,
		},
	}

	return &resp, nil
}

func (h *bookPageHandler) GetBookPagesByOffset(ctx context.Context, input *schema.GetBookPagesOffsetInput) (*schema.GetBookPagesOutput, error) {
	bookPages, err := h.BookPageUseCase.GetBookPagesByOffset(ctx, input.BookID, input.CenterPage, input.Offset)
	if err != nil {
		return nil, err
	}

	bookPageOutputs := convertBookPages(bookPages)

	resp := schema.GetBookPagesOutput{
		Body: schema.GetBookPagesOutputBody{
			Data: bookPageOutputs,
		},
	}

	return &resp, nil
}

func (h *bookPageHandler) GetBookPageByID(ctx context.Context, input *schema.GetBookPageByIDInput) (*schema.GetBookPageByIDOutput, error) {
	bookPage, err := h.BookPageUseCase.GetBookPageByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	bookPageOutput := convertBookPage(bookPage)

	resp := schema.GetBookPageByIDOutput{
		Body: bookPageOutput,
	}

	return &resp, nil
}

func convertBookPages(bookPages []entity.BookPage) []schema.BookPageOutput {
	bookPageOutputs := make([]schema.BookPageOutput, len(bookPages))
	for i, r := range bookPages {
		bookPageOutputs[i] = convertBookPage(r)
	}
	return bookPageOutputs
}

func convertBookPage(bookPage entity.BookPage) schema.BookPageOutput {
	return schema.BookPageOutput{
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
