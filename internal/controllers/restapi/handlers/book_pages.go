package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"

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
	BooksRepository     mongo_repositories.BooksRepository
	BookPagesRepository mongo_repositories.BookPagesRepository
}

func NewBookPagesHandler(booksRepo mongo_repositories.BooksRepository, bookPagesRepo mongo_repositories.BookPagesRepository) *bookPagesHandler {
	return &bookPagesHandler{
		BooksRepository:     booksRepo,
		BookPagesRepository: bookPagesRepo,
	}
}

func (h *bookPagesHandler) CreateBookPage(ctx context.Context, input *models.CreateBookPageInput) (*models.CreateBookPageOutput, error) {
	bookOID, err := bson.ObjectIDFromHex(input.Body.BookID)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %w", err)
	}
	_, err = h.BooksRepository.GetBookByID(ctx, input.Body.BookID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book info: %w", err)
	}

	currentTime := time.Now()

	metadata := input.Body.Metadata
	record := &entities.BookPage{
		BookID:     bookOID,
		PageNumber: input.Body.PageNumber,
		Content:    input.Body.Content,
		Metadata: entities.BookPageMetadata{
			IsBookmarked: metadata.IsBookmarked,
			Highlight:    metadata.Highlight,
		},
		AttachedImageFileID: input.Body.AttachedImageFileID,
		CreatedAt:           currentTime,
		ModifiedAt:          currentTime,
	}

	id, err := h.BookPagesRepository.CreateBookPage(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to create book page: %w", err)
	}

	resp := &models.CreateBookPageOutput{
		Body: models.CreateBookPageOutputBody{
			ID: id,
		},
	}

	return resp, nil
}

func (h *bookPagesHandler) GetBookPages(ctx context.Context, input *models.GetBookPagesInput) (*models.GetBookPagesOutput, error) {
	var records []entities.BookPage
	var err error

	if input.GetAll {
		records, err = h.BookPagesRepository.GetBookPagesByBookID(ctx, input.BookID)
	} else {
		records, err = h.BookPagesRepository.GetBookpagesByBookIDWithPagination(ctx, input.BookID, input.PageSize, input.PageNumber)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages: %w", err)
	}

	results := convertBookPages(records)

	resp := &models.GetBookPagesOutput{
		Body: models.GetBookPagesOutputBody{
			Data: results,
		},
	}

	return resp, nil
}

func (h *bookPagesHandler) GetBookPagesByRange(ctx context.Context, input *models.GetBookPagesRangeInput) (*models.GetBookPagesOutput, error) {
	records, err := h.BookPagesRepository.GetBookpagesByPageRange(ctx, input.BookID, input.StartPage, input.EndPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by range: %w", err)
	}

	results := convertBookPages(records)

	resp := &models.GetBookPagesOutput{
		Body: models.GetBookPagesOutputBody{
			Data: results,
		},
	}

	return resp, nil
}

func (h *bookPagesHandler) GetBookPagesByOffset(ctx context.Context, input *models.GetBookPagesOffsetInput) (*models.GetBookPagesOutput, error) {
	records, err := h.BookPagesRepository.GetBookpagesAroundPageNumber(ctx, input.BookID, input.CenterPage, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book pages by offset: %w", err)
	}

	results := convertBookPages(records)

	resp := &models.GetBookPagesOutput{
		Body: models.GetBookPagesOutputBody{
			Data: results,
		},
	}

	return resp, nil
}

func (h *bookPagesHandler) GetBookPageByID(ctx context.Context, input *models.GetBookPageByIDInput) (*models.GetBookPageByIDOutput, error) {
	record, err := h.BookPagesRepository.GetBookPageByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book page: %w", err)
	}

	output := convertBookPage(*record)

	resp := &models.GetBookPageByIDOutput{
		Body: output,
	}

	return resp, nil
}

func convertBookPages(records []entities.BookPage) []models.BookPageOutput {
	results := make([]models.BookPageOutput, len(records))
	for i, r := range records {
		results[i] = convertBookPage(r)
	}
	return results
}

func convertBookPage(record entities.BookPage) models.BookPageOutput {
	metadata := record.Metadata

	return models.BookPageOutput{
		ID:         record.ID.Hex(),
		BookID:     record.BookID.Hex(),
		PageNumber: record.PageNumber,
		Content:    record.Content,
		Metadata: models.BookPageMetadata{
			IsBookmarked: metadata.IsBookmarked,
			Highlight:    metadata.Highlight,
		},
		AttachedImageFileID: record.AttachedImageFileID,
	}
}
