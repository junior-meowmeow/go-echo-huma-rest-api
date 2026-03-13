package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterBookPageGroup(api huma.API, h *handler.Handlers) {
	bookPageGroup := huma.NewGroup(api, "/book_pages")

	RegisterBookPageRoutes(bookPageGroup, h)
}

func RegisterBookPageRoutes(api huma.API, h *handler.Handlers) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-book-page",
		Method:        http.MethodPost,
		Path:          "/",
		Summary:       "Create Book Page",
		Description:   "Create a new book page.",
		Tags:          []string{"Book Pages"},
		DefaultStatus: http.StatusCreated,
	}, h.BookPage.CreateBookPage)

	huma.Register(api, huma.Operation{
		OperationID: "get-book-pages",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Get Book Pages",
		Description: "Get book pages in a book.",
		Tags:        []string{"Book Pages"},
	}, h.BookPage.GetBookPages)

	huma.Register(api, huma.Operation{
		OperationID: "get-book-pages-range",
		Method:      http.MethodGet,
		Path:        "/range",
		Summary:     "Get Book Pages (Page Range)",
		Description: "Get book pages within start/end page number.",
		Tags:        []string{"Book Pages"},
	}, h.BookPage.GetBookPagesByRange)

	huma.Register(api, huma.Operation{
		OperationID: "get-book-pages-offset",
		Method:      http.MethodGet,
		Path:        "/offset",
		Summary:     "Get Book Pages (Offset)",
		Description: "Get N book pages before and after specified page number.",
		Tags:        []string{"Book Pages"},
	}, h.BookPage.GetBookPagesByOffset)

	huma.Register(api, huma.Operation{
		OperationID: "get-book-page-by-id",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Get Book Page",
		Description: "Get a book page by ID.",
		Tags:        []string{"Book Pages"},
	}, h.BookPage.GetBookPageByID)
}
