package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type BookHandler interface {
	CreateBook(ctx context.Context, input *schema.CreateBookInput) (*schema.CreateBookOutput, error)
	GetBooks(ctx context.Context, input *schema.GetBooksInput) (*schema.GetBooksOutput, error)
	GetBookByID(ctx context.Context, input *schema.GetBookByIDInput) (*schema.GetBookByIDOutput, error)
}

type bookHandler struct {
	BookUseCase usecase.BookUseCase
}

func NewBookHandler(bookUseCase usecase.BookUseCase) *bookHandler {
	return &bookHandler{
		BookUseCase: bookUseCase,
	}
}

func (h *bookHandler) CreateBook(ctx context.Context, input *schema.CreateBookInput) (*schema.CreateBookOutput, error) {
	book := &entity.Book{
		Name:        input.Body.Name,
		Description: input.Body.Description,
		Metadata: entity.BookMetadata{
			Author: input.Body.Metadata.Author,
			ISBN:   input.Body.Metadata.ISBN,
			Genre:  input.Body.Metadata.Genre,
		},
		CoverImageFileID: input.Body.CoverImageFileID,
	}

	id, err := h.BookUseCase.CreateBook(ctx, book)
	if err != nil {
		return nil, err
	}

	resp := schema.CreateBookOutput{
		Body: schema.CreateBookOutputBody{
			ID: id,
		},
	}

	return &resp, nil
}

func (h *bookHandler) GetBooks(ctx context.Context, input *schema.GetBooksInput) (*schema.GetBooksOutput, error) {
	var books []entity.Book
	var err error

	if input.GetAll {
		books, err = h.BookUseCase.GetAllBooks(ctx)
	} else {
		books, err = h.BookUseCase.GetBooksWithPagination(ctx, input.PageSize, input.PageNumber)
	}

	if err != nil {
		return nil, err
	}

	bookOutputs := convertBooks(books)

	resp := schema.GetBooksOutput{
		Body: schema.GetBooksOutputBody{
			Data: bookOutputs,
		},
	}

	return &resp, nil
}

func (h *bookHandler) GetBookByID(ctx context.Context, input *schema.GetBookByIDInput) (*schema.GetBookByIDOutput, error) {
	book, err := h.BookUseCase.GetBookByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	bookOutput := convertBook(book)

	resp := schema.GetBookByIDOutput{
		Body: bookOutput,
	}

	return &resp, nil
}

func convertBooks(books []entity.Book) []schema.BookOutput {
	bookOutputs := make([]schema.BookOutput, len(books))
	for i, r := range books {
		bookOutputs[i] = convertBook(r)
	}
	return bookOutputs
}

func convertBook(book entity.Book) schema.BookOutput {
	return schema.BookOutput{
		ID:          book.ID,
		Name:        book.Name,
		Description: book.Description,
		Metadata: schema.BookMetadata{
			Author: book.Metadata.Author,
			ISBN:   book.Metadata.ISBN,
			Genre:  book.Metadata.Genre,
		},
		CoverImageFileID: book.CoverImageFileID,
		CreatedAt:        book.CreatedAt,
	}
}
