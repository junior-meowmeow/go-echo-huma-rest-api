package schema

import (
	"time"
)

type Review struct {
	ID        string    `json:"id" doc:"Review ID"`
	Author    string    `json:"author" doc:"Name of the review author"`
	Rating    int       `json:"rating" doc:"Rating from 1 to 5"`
	Message   string    `json:"message,omitempty" doc:"Review message"`
	CreatedAt time.Time `json:"createdAt" doc:"Timestamp when the review was created"`
}

type CreateReviewRequest struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Name of the review author" example:"User"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5" example:"2"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message" example:"Lorem Ipsum"`
	}
}

type CreateReviewResponse struct{}

type GetReviewsRequest struct{}

type GetReviewsResponse struct {
	Body struct {
		Data []Review `json:"data" doc:"List of reviews"`
	}
}
