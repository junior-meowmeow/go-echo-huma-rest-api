package schema

import (
	"time"
)

type BookMetadata struct {
	Author string `json:"author" doc:"Author name"`
	ISBN   string `json:"isbn" doc:"ISBN of the book"`
	Genre  string `json:"genre" doc:"Book genre(s)"`
}

type CreateBookInput struct {
	Body struct {
		Name             string       `json:"name" required:"true" maxLength:"100" doc:"books name" example:"New Book"`
		Description      string       `json:"description,omitempty" maxLength:"500" doc:"Description"`
		Metadata         BookMetadata `json:"metadata" doc:"Metadata of the book"`
		CoverImageFileID string       `json:"coverImageFileId,omitempty" doc:"File ID for preview image"`
	}
}

type CreateBookOutputBody struct {
	ID string `json:"id" doc:"Created Book ID"`
}

type CreateBookOutput struct {
	Body CreateBookOutputBody
}

type BookOutput struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Description      string       `json:"description,omitempty"`
	Metadata         BookMetadata `json:"metadata"`
	CoverImageFileID string       `json:"coverImageFileId,omitempty"`
	CreatedAt        time.Time    `json:"createdAt"`
}

type GetBooksInput struct {
	GetAll     bool  `query:"all" required:"true" default:"false" doc:"If true, returns all books ignoring pagination"`
	PageNumber int64 `query:"pageNumber" default:"1" minimum:"1" doc:"Page number"`
	PageSize   int64 `query:"pageSize" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
}

type GetBooksOutputBody struct {
	Data []BookOutput `json:"data"`
}

type GetBooksOutput struct {
	Body GetBooksOutputBody
}

type GetBookByIDInput struct {
	ID string `path:"id" pattern:"^[a-fA-F0-9]{24}$" required:"true" doc:"Book ID"`
}

type GetBookByIDOutput struct {
	Body BookOutput
}
