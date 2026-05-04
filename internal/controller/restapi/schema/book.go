package schema

import (
	"time"
)

type Book struct {
	ID               string       `json:"id" doc:"Book ID" readOnly:"true"`
	Name             string       `json:"name" maxLength:"100" doc:"Book name" example:"New Book"`
	Description      string       `json:"description" maxLength:"500" doc:"Book description"`
	Metadata         BookMetadata `json:"metadata" doc:"Metadata of the book"`
	CoverImageFileID string       `json:"coverImageFileId,omitempty" doc:"File ID of the cover image"`
	CreatedAt        time.Time    `json:"createdAt" doc:"Timestamp when the book was created" readOnly:"true"`
}

type BookMetadata struct {
	Author string `json:"author" doc:"Author name"`
	ISBN   string `json:"isbn" doc:"ISBN of the book"`
	Genre  string `json:"genre,omitempty" doc:"Book genre(s)"`
}

type CreateBookRequest struct {
	Body Book
}

type CreateBookResponse struct {
	Body struct {
		ID string `json:"id" doc:"Created Book ID"`
	}
}

type GetBooksRequest struct {
	GetAll     bool  `query:"all" required:"true" default:"false" doc:"If true, returns all items ignoring pagination"`
	PageNumber int64 `query:"pageNumber" minimum:"1" default:"1" doc:"Page number"`
	PageSize   int64 `query:"pageSize" minimum:"1" maximum:"100" default:"20" doc:"Items per page"`
}

type GetBooksResponse struct {
	Body struct {
		Data []Book `json:"data" doc:"List of books"`
	}
}

type GetBookByIDRequest struct {
	ID string `path:"id" required:"true" format:"uuid" doc:"Book ID"`
}

type GetBookByIDResponse struct {
	Body Book
}
