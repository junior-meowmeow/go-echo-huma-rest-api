package models

import "time"

type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Author of the review" example:"User"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5" example:"2"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message" example:"Lorem Ipsum"`
	}
}

type ReviewOutput struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetReviewsOutput struct {
	Body []ReviewOutput `json:"body"`
}
