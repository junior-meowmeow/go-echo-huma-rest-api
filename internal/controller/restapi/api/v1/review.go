package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
	customMiddleware "github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/middleware"
)

func RegisterReviewGroup(public huma.API, protected huma.API, h *handler.Handlers) {
	publicGroup := huma.NewGroup(public, "/reviews")
	protectedGroup := huma.NewGroup(protected, "/reviews")

	RegisterReviewRoutes(publicGroup, protectedGroup, h)
}

func RegisterReviewRoutes(public huma.API, protected huma.API, h *handler.Handlers) {
	huma.Register(protected, huma.Operation{
		OperationID:   "create-review",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Post new review",
		Description:   "Post a new review to database.",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusCreated,
	}, h.Review.CreateReview)

	huma.Register(protected, huma.Operation{
		OperationID: "get-reviews",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "Get all reviews",
		Description: "Get all reviews from database.",
		Tags:        []string{"Reviews"},
		Middlewares: huma.Middlewares{
			customMiddleware.RequireAdminRole(protected),
		},
	}, h.Review.GetReviews)
}
