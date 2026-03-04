package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository"
)

func (h *Handler) CreateBook(ctx context.Context, input *models.CreateBookInput) (*models.CreateBookOutput, error) {
	currentTime := time.Now()

	record := &repository.BookRecord{
		Name:        input.Body.Name,
		Description: input.Body.Description,
		Metadata: repository.BookMetadata{
			Author: input.Body.Metadata.Author,
			ISBN:   input.Body.Metadata.ISBN,
			Genre:  input.Body.Metadata.Genre,
		},
		CoverImageFileID: input.Body.CoverImageFileID,
		CreatedAt:        currentTime,
		ModifiedAt:       currentTime,
	}

	id, err := h.Books.CreateBook(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	resp := &models.CreateBookOutput{
		Body: models.CreateBookOutputBody{
			ID: id,
		},
	}

	return resp, nil

}

func (h *Handler) GetBooks(ctx context.Context, input *models.GetBooksInput) (*models.GetBooksOutput, error) {
	var records []repository.BookRecord
	var err error

	if input.GetAll {
		records, err = h.Books.GetAllBooks(ctx)
	} else {
		records, err = h.Books.GetBooksWithPagination(ctx, input.PageSize, input.PageNumber)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch books: %w", err)
	}

	results := make([]models.BookOutput, len(records))
	for i, r := range records {
		results[i] = models.BookOutput{
			ID:          r.ID.Hex(),
			Name:        r.Name,
			Description: r.Description,
			Metadata: models.BookMetadata{
				Author: r.Metadata.Author,
				ISBN:   r.Metadata.ISBN,
				Genre:  r.Metadata.Genre,
			},
			CoverImageFileID: r.CoverImageFileID,
			CreatedAt:        r.CreatedAt,
		}
	}

	resp := &models.GetBooksOutput{
		Body: models.GetBooksOutputBody{
			Data: results,
		},
	}

	return resp, nil
}

func (h *Handler) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.GetBookByIDOutput, error) {
	record, err := h.Books.GetBookByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book: %w", err)
	}

	resp := &models.GetBookByIDOutput{
		Body: models.BookOutput{
			ID:          record.ID.Hex(),
			Name:        record.Name,
			Description: record.Description,
			Metadata: models.BookMetadata{
				Author: record.Metadata.Author,
				ISBN:   record.Metadata.ISBN,
				Genre:  record.Metadata.Genre,
			},
			CoverImageFileID: record.CoverImageFileID,
			CreatedAt:        record.CreatedAt,
		},
	}

	return resp, nil
}
