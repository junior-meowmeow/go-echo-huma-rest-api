package v1

import (
	"net/http"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGroup(api huma.API, h *handlers.Handlers) {
	v1 := huma.NewGroup(api, "/v1")

	v1.UseSimpleModifier(func(op *huma.Operation) {
		op.OperationID = op.OperationID + "-v1"
		op.Summary = op.Summary + " (V1)"
	})

	RegisterRoutes(v1, h)
}

func RegisterRoutes(api huma.API, h *handlers.Handlers) {
	// POST /reviews
	huma.Register(api, huma.Operation{
		OperationID:   "post-review",
		Method:        http.MethodPost,
		Path:          "/reviews",
		Summary:       "Post new review",
		Description:   "Post a new review to database.",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusCreated,
	}, h.Reviews.PostReview)

	// GET /reviews
	huma.Register(api, huma.Operation{
		OperationID: "get-reviews",
		Method:      http.MethodGet,
		Path:        "/reviews",
		Summary:     "Get all reviews",
		Description: "Get all reviews from database.",
		Tags:        []string{"Reviews"},
	}, h.Reviews.GetReviews)

	// POST /files/upload
	huma.Register(api, huma.Operation{
		OperationID: "upload-s3-file",
		Method:      http.MethodPost,
		Path:        "/files/upload",
		Summary:     "Upload file to object storage",
		Description: "Upload a file to object storage.",
		Tags:        []string{"Files"},
	}, h.Files.UploadFile)

	// GET /files/download/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-s3-file",
		Method:      http.MethodGet,
		Path:        "/files/download/{id}",
		Summary:     "Get file from object storage",
		Description: "Get a file from object storage.",
		Tags:        []string{"Files"},
	}, h.Files.GetFileDownloadLink)

	// POST /books
	huma.Register(api, huma.Operation{
		OperationID:   "create-book",
		Method:        http.MethodPost,
		Path:          "/books",
		Summary:       "Create Book",
		Description:   "Create a new book.",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusCreated,
	}, h.Books.CreateBook)

	// GET /books
	huma.Register(api, huma.Operation{
		OperationID: "list-books",
		Method:      http.MethodGet,
		Path:        "/books",
		Summary:     "List Books",
		Description: "Get a list of books.",
		Tags:        []string{"Books"},
	}, h.Books.GetBooks)

	// GET /books/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-book-by-id",
		Method:      http.MethodGet,
		Path:        "/books/{id}",
		Summary:     "Get Book",
		Description: "Get a book by ID.",
		Tags:        []string{"Books"},
	}, h.Books.GetBookByID)

	// POST /book_pages
	huma.Register(api, huma.Operation{
		OperationID:   "create-book-page",
		Method:        http.MethodPost,
		Path:          "/book_pages",
		Summary:       "Create Book Page",
		Description:   "Create a new book page.",
		Tags:          []string{"Book Pages"},
		DefaultStatus: http.StatusCreated,
	}, h.BookPages.CreateBookPage)

	// GET /book_pages
	huma.Register(api, huma.Operation{
		OperationID: "list-book-pages",
		Method:      http.MethodGet,
		Path:        "/book_pages",
		Summary:     "List Book Pages",
		Description: "List book pages in a book.",
		Tags:        []string{"Book Pages"},
	}, h.BookPages.GetBookPages)

	// GET /book_pages/range
	huma.Register(api, huma.Operation{
		OperationID: "list-book-pages-range",
		Method:      http.MethodGet,
		Path:        "/book_pages/range",
		Summary:     "List Book Pages (Time Range)",
		Description: "Get book pages within start/end page number.",
		Tags:        []string{"Book Pages"},
	}, h.BookPages.GetBookPagesByRange)

	// GET /book_pages/offset
	huma.Register(api, huma.Operation{
		OperationID: "list-book-pages-offset",
		Method:      http.MethodGet,
		Path:        "/book_pages/offset",
		Summary:     "List Book Pages (Offset)",
		Description: "Get N book pages before and after specified page number.",
		Tags:        []string{"Book Pages"},
	}, h.BookPages.GetBookPagesByOffset)

	// GET /book_pages/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-book-page-by-id",
		Method:      http.MethodGet,
		Path:        "/book_pages/{id}",
		Summary:     "Get Book Page",
		Description: "Get a book page by ID.",
		Tags:        []string{"Book Pages"},
	}, h.BookPages.GetBookPageByID)
}
