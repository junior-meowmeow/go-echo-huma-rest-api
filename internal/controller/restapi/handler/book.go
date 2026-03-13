package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type BookHandler interface {
	CreateBook(ctx context.Context, request *schema.CreateBookRequest) (*schema.CreateBookResponse, error)
	GetBooks(ctx context.Context, request *schema.GetBooksRequest) (*schema.GetBooksResponse, error)
	GetBookByID(ctx context.Context, request *schema.GetBookByIDRequest) (*schema.GetBookByIDResponse, error)
}

type bookHandler struct {
	BookUseCase usecase.BookUseCase
}

func NewBookHandler(bookUseCase usecase.BookUseCase) *bookHandler {
	return &bookHandler{
		BookUseCase: bookUseCase,
	}
}

func (h *bookHandler) CreateBook(ctx context.Context, request *schema.CreateBookRequest) (*schema.CreateBookResponse, error) {
	book := &entity.Book{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Metadata: entity.BookMetadata{
			Author: request.Body.Metadata.Author,
			ISBN:   request.Body.Metadata.ISBN,
			Genre:  request.Body.Metadata.Genre,
		},
		CoverImageFileID: request.Body.CoverImageFileID,
	}

	id, err := h.BookUseCase.CreateBook(ctx, book)
	if err != nil {
		return nil, err
	}

	resp := schema.CreateBookResponse{}
	resp.Body.ID = id

	return &resp, nil
}

func (h *bookHandler) GetBooks(ctx context.Context, request *schema.GetBooksRequest) (*schema.GetBooksResponse, error) {
	var books []entity.Book
	var err error

	if request.GetAll {
		books, err = h.BookUseCase.GetAllBooks(ctx)
	} else {
		books, err = h.BookUseCase.GetBooksWithPagination(ctx, request.PageSize, request.PageNumber)
	}

	if err != nil {
		return nil, err
	}

	bookOutputs := convertBooks(books)

	resp := schema.GetBooksResponse{}
	resp.Body.Data = bookOutputs

	return &resp, nil
}

func (h *bookHandler) GetBookByID(ctx context.Context, request *schema.GetBookByIDRequest) (*schema.GetBookByIDResponse, error) {
	book, err := h.BookUseCase.GetBookByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	bookOutput := convertBook(book)

	resp := schema.GetBookByIDResponse{}
	resp.Body = bookOutput

	return &resp, nil
}

func convertBooks(books []entity.Book) []schema.Book {
	bookOutputs := make([]schema.Book, len(books))
	for i, r := range books {
		bookOutputs[i] = convertBook(r)
	}
	return bookOutputs
}

func convertBook(book entity.Book) schema.Book {
	return schema.Book{
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
