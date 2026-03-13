package schema

type BookPage struct {
	ID                  string           `json:"id" doc:"Book Page ID"`
	BookID              string           `json:"bookId" doc:"Parent Book ID"`
	PageNumber          int64            `json:"pageNumber" doc:"Page number within the book"`
	Content             string           `json:"content" doc:"Text content of the page"`
	Metadata            BookPageMetadata `json:"metadata" doc:"Metadata of the book page"`
	AttachedImageFileID string           `json:"attachedImageFileId,omitempty" doc:"File ID of the attached image"`
}

type BookPageMetadata struct {
	IsBookmarked bool   `json:"isBookmarked" doc:"Is the page bookmarked or not"`
	Highlight    string `json:"highlight" doc:"Highlight content"`
}

type ParentBookID struct {
	BookID string `query:"bookId" pattern:"^[a-fA-F0-9]{24}$" required:"true" doc:"Parent Book ID"`
}

type CreateBookPageRequest struct {
	Body struct {
		BookID              string           `json:"bookID" pattern:"^[a-fA-F0-9]{24}$" doc:"Parent Book ID"`
		PageNumber          int64            `json:"pageNumber" doc:"Page number within the book"`
		Content             string           `json:"content" doc:"Text content of the page"`
		Metadata            BookPageMetadata `json:"metadata" doc:"Metadata of the book page"`
		AttachedImageFileID string           `json:"attachedImageFileId,omitempty" doc:"File ID of the attached image"`
	}
}

type CreateBookPageResponse struct {
	Body struct {
		ID string `json:"id"`
	}
}

type GetBookPagesRequest struct {
	ParentBookID
	GetAll     bool  `query:"all" required:"true" default:"false" doc:"If true, returns all items ignoring pagination"`
	PageNumber int64 `query:"pageNumber" default:"1" minimum:"1" doc:"Page number"`
	PageSize   int64 `query:"pageSize" default:"50" minimum:"1" maximum:"500" doc:"Items per page"`
}

type GetBookPagesResponse struct {
	Body struct {
		Data []BookPage `json:"data" doc:"List of book pages"`
	}
}

type GetBookPagesRangeRequest struct {
	ParentBookID
	StartPage int64 `query:"startPage" required:"true" doc:"Start page number (inclusive)"`
	EndPage   int64 `query:"endPage" required:"true" doc:"End page number (inclusive)"`
}

type GetBookPagesOffsetRequest struct {
	ParentBookID
	CenterPage int64 `query:"centerPage" required:"true" doc:"Center page used as the reference point"`
	Offset     int64 `query:"offset" required:"true" default:"10" doc:"Number of pages to include before and after the center page"`
}

type GetBookPageByIDRequest struct {
	ID string `path:"id" pattern:"^[a-fA-F0-9]{24}$" doc:"Book Page ID"`
}

type GetBookPageByIDResponse struct {
	Body BookPage
}
