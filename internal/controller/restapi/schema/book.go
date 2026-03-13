package schema

import (
	"time"
)

type Book struct {
	ID               string       `json:"id" doc:"Book ID"`
	Name             string       `json:"name" doc:"Book name"`
	Description      string       `json:"description" doc:"Book description"`
	Metadata         BookMetadata `json:"metadata" doc:"Metadata of the book"`
	CoverImageFileID string       `json:"coverImageFileId,omitempty" doc:"File ID of the cover image"`
	CreatedAt        time.Time    `json:"createdAt" doc:"Timestamp when the book was created"`
}

type BookMetadata struct {
	Author string `json:"author" doc:"Author name"`
	ISBN   string `json:"isbn" doc:"ISBN of the book"`
	Genre  string `json:"genre" doc:"Book genre(s)"`
}

type CreateBookRequest struct {
	Body struct {
		Name             string       `json:"name" required:"true" maxLength:"100" doc:"Book name" example:"New Book"`
		Description      string       `json:"description,omitempty" maxLength:"500" doc:"Book description"`
		Metadata         BookMetadata `json:"metadata" doc:"Metadata of the book"`
		CoverImageFileID string       `json:"coverImageFileId,omitempty" doc:"File ID of the book cover image"`
	}
}

type CreateBookResponse struct {
	Body struct {
		ID string `json:"id" doc:"Created Book ID"`
	}
}

type GetBooksRequest struct {
	GetAll     bool  `query:"all" required:"true" default:"false" doc:"If true, returns all items ignoring pagination"`
	PageNumber int64 `query:"pageNumber" default:"1" minimum:"1" doc:"Page number"`
	PageSize   int64 `query:"pageSize" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
}

type GetBooksResponse struct {
	Body struct {
		Data []Book `json:"data"`
	}
}

type GetBookByIDRequest struct {
	ID string `path:"id" json:"id" pattern:"^[a-fA-F0-9]{24}$" required:"true" doc:"Book ID"`
}

type GetBookByIDResponse struct {
	Body Book
}
