package v1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/handler"
)

func RegisterReviewGroup(api huma.API, h *handler.Handlers) {
	reviewGroup := huma.NewGroup(api, "/reviews")

	RegisterReviewRoutes(reviewGroup, h)
}

func RegisterReviewRoutes(api huma.API, h *handler.Handlers) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-review",
		Method:        http.MethodPost,
		Path:          "",
		Summary:       "Post new review",
		Description:   "Post a new review to database.",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusCreated,
	}, h.Review.CreateReview)

	huma.Register(api, huma.Operation{
		OperationID: "get-reviews",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "Get all reviews",
		Description: "Get all reviews from database.",
		Tags:        []string{"Reviews"},
	}, h.Review.GetReviews)
}
