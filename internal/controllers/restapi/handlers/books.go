package handlers

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"
)

type BooksHandler interface {
	CreateBook(ctx context.Context, input *models.CreateBookInput) (*models.CreateBookOutput, error)
	GetBooks(ctx context.Context, input *models.GetBooksInput) (*models.GetBooksOutput, error)
	GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.GetBookByIDOutput, error)
}

type booksHandler struct {
	BooksUseCase usecases.BooksUseCase
}

func NewBooksHandler(booksUseCase usecases.BooksUseCase) *booksHandler {
	return &booksHandler{
		BooksUseCase: booksUseCase,
	}
}

func (h *booksHandler) CreateBook(ctx context.Context, input *models.CreateBookInput) (*models.CreateBookOutput, error) {
	book := &entities.Book{
		Name:        input.Body.Name,
		Description: input.Body.Description,
		Metadata: entities.BookMetadata{
			Author: input.Body.Metadata.Author,
			ISBN:   input.Body.Metadata.ISBN,
			Genre:  input.Body.Metadata.Genre,
		},
		CoverImageFileID: input.Body.CoverImageFileID,
	}

	id, err := h.BooksUseCase.CreateBook(ctx, book)
	if err != nil {
		return nil, err
	}

	resp := models.CreateBookOutput{
		Body: models.CreateBookOutputBody{
			ID: id,
		},
	}

	return &resp, nil
}

func (h *booksHandler) GetBooks(ctx context.Context, input *models.GetBooksInput) (*models.GetBooksOutput, error) {
	var books []entities.Book
	var err error

	if input.GetAll {
		books, err = h.BooksUseCase.GetAllBooks(ctx)
	} else {
		books, err = h.BooksUseCase.GetBooksWithPagination(ctx, input.PageSize, input.PageNumber)
	}

	if err != nil {
		return nil, err
	}

	bookOutputs := convertBooks(books)

	resp := models.GetBooksOutput{
		Body: models.GetBooksOutputBody{
			Data: bookOutputs,
		},
	}

	return &resp, nil
}

func (h *booksHandler) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.GetBookByIDOutput, error) {
	book, err := h.BooksUseCase.GetBookByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	bookOutput := convertBook(book)

	resp := models.GetBookByIDOutput{
		Body: bookOutput,
	}

	return &resp, nil
}

func convertBooks(books []entities.Book) []models.BookOutput {
	bookOutputs := make([]models.BookOutput, len(books))
	for i, r := range books {
		bookOutputs[i] = convertBook(r)
	}
	return bookOutputs
}

func convertBook(book entities.Book) models.BookOutput {
	return models.BookOutput{
		ID:          book.ID,
		Name:        book.Name,
		Description: book.Description,
		Metadata: models.BookMetadata{
			Author: book.Metadata.Author,
			ISBN:   book.Metadata.ISBN,
			Genre:  book.Metadata.Genre,
		},
		CoverImageFileID: book.CoverImageFileID,
		CreatedAt:        book.CreatedAt,
	}
}
