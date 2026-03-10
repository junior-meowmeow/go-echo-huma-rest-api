package handlers

import (
	"context"
	"fmt"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookPagesHandler interface {
	CreateBookPage(ctx context.Context, input *models.CreateBookPageInput) (*models.CreateBookPageOutput, error)
	GetBookPages(ctx context.Context, input *models.GetBookPagesInput) (*models.GetBookPagesOutput, error)
	GetBookPagesByRange(ctx context.Context, input *models.GetBookPagesRangeInput) (*models.GetBookPagesOutput, error)
	GetBookPagesByOffset(ctx context.Context, input *models.GetBookPagesOffsetInput) (*models.GetBookPagesOutput, error)
	GetBookPageByID(ctx context.Context, input *models.GetBookPageByIDInput) (*models.GetBookPageByIDOutput, error)
}

type bookPagesHandler struct {
	BookPagesUseCase usecases.BookPagesUseCase
}

func NewBookPagesHandler(bookPagesUseCase usecases.BookPagesUseCase) *bookPagesHandler {
	return &bookPagesHandler{
		BookPagesUseCase: bookPagesUseCase,
	}
}

func (h *bookPagesHandler) CreateBookPage(ctx context.Context, input *models.CreateBookPageInput) (*models.CreateBookPageOutput, error) {
	bookOID, err := bson.ObjectIDFromHex(input.Body.BookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %w", err)
	}

	metadata := input.Body.Metadata
	bookPage := &entities.BookPage{
		BookID:     bookOID,
		PageNumber: input.Body.PageNumber,
		Content:    input.Body.Content,
		Metadata: entities.BookPageMetadata{
			IsBookmarked: metadata.IsBookmarked,
			Highlight:    metadata.Highlight,
		},
		AttachedImageFileID: input.Body.AttachedImageFileID,
	}

	id, err := h.BookPagesUseCase.CreateBookPage(ctx, bookPage)
	if err != nil {
		return nil, err
	}

	resp := models.CreateBookPageOutput{
		Body: models.CreateBookPageOutputBody{
			ID: id,
		},
	}

	return &resp, nil
}

func (h *bookPagesHandler) GetBookPages(ctx context.Context, input *models.GetBookPagesInput) (*models.GetBookPagesOutput, error) {
	var bookPages []entities.BookPage
	var err error

	if input.GetAll {
		bookPages, err = h.BookPagesUseCase.GetAllBookPages(ctx, input.BookID)
	} else {
		bookPages, err = h.BookPagesUseCase.GetBookPagesWithPagination(ctx, input.BookID, input.PageSize, input.PageNumber)
	}

	if err != nil {
		return nil, err
	}

	bookPagesOutput := convertBookPages(bookPages)

	resp := models.GetBookPagesOutput{
		Body: models.GetBookPagesOutputBody{
			Data: bookPagesOutput,
		},
	}

	return &resp, nil
}

func (h *bookPagesHandler) GetBookPagesByRange(ctx context.Context, input *models.GetBookPagesRangeInput) (*models.GetBookPagesOutput, error) {
	bookPages, err := h.BookPagesUseCase.GetBookPagesByRange(ctx, input.BookID, input.StartPage, input.EndPage)
	if err != nil {
		return nil, err
	}

	bookPageOutputs := convertBookPages(bookPages)

	resp := models.GetBookPagesOutput{
		Body: models.GetBookPagesOutputBody{
			Data: bookPageOutputs,
		},
	}

	return &resp, nil
}

func (h *bookPagesHandler) GetBookPagesByOffset(ctx context.Context, input *models.GetBookPagesOffsetInput) (*models.GetBookPagesOutput, error) {
	bookPages, err := h.BookPagesUseCase.GetBookPagesByOffset(ctx, input.BookID, input.CenterPage, input.Offset)
	if err != nil {
		return nil, err
	}

	bookPageOutputs := convertBookPages(bookPages)

	resp := models.GetBookPagesOutput{
		Body: models.GetBookPagesOutputBody{
			Data: bookPageOutputs,
		},
	}

	return &resp, nil
}

func (h *bookPagesHandler) GetBookPageByID(ctx context.Context, input *models.GetBookPageByIDInput) (*models.GetBookPageByIDOutput, error) {
	bookPage, err := h.BookPagesUseCase.GetBookPageByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	bookPageOutput := convertBookPage(bookPage)

	resp := models.GetBookPageByIDOutput{
		Body: bookPageOutput,
	}

	return &resp, nil
}

func convertBookPages(bookPages []entities.BookPage) []models.BookPageOutput {
	bookPageOutputs := make([]models.BookPageOutput, len(bookPages))
	for i, r := range bookPages {
		bookPageOutputs[i] = convertBookPage(r)
	}
	return bookPageOutputs
}

func convertBookPage(bookPage entities.BookPage) models.BookPageOutput {
	return models.BookPageOutput{
		ID:         bookPage.ID.Hex(),
		BookID:     bookPage.BookID.Hex(),
		PageNumber: bookPage.PageNumber,
		Content:    bookPage.Content,
		Metadata: models.BookPageMetadata{
			IsBookmarked: bookPage.Metadata.IsBookmarked,
			Highlight:    bookPage.Metadata.Highlight,
		},
		AttachedImageFileID: bookPage.AttachedImageFileID,
	}
}
