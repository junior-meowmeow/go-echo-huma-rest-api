package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterBookGroup(api huma.API, h *handler.Handlers) {
	bookGroup := huma.NewGroup(api, "/books")

	RegisterBookRoutes(bookGroup, h)
}

func RegisterBookRoutes(api huma.API, h *handler.Handlers) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-book",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Create Book",
		Description:   "Create a new book.",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusCreated,
	}, h.Book.CreateBook)

	huma.Register(api, huma.Operation{
		OperationID: "get-books",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "Get Books",
		Description: "Get a list of books.",
		Tags:        []string{"Books"},
	}, h.Book.GetBooks)

	huma.Register(api, huma.Operation{
		OperationID: "get-book-by-id",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Get Book",
		Description: "Get a book by ID.",
		Tags:        []string{"Books"},
	}, h.Book.GetBookByID)
}
