package models

type BookPageMetadata struct {
	IsBookmarked bool   `json:"isBookmarked" doc:"Is the page bookmarked or not"`
	Highlight    string `json:"highlight" doc:"Highlight content"`
}

type CreateBookPageInput struct {
	Body struct {
		BookID              string           `json:"bookID" pattern:"^[a-fA-F0-9]{24}$" doc:"Parent Book ID"`
		PageNumber          int64            `json:"pageNumber" doc:"Page number"`
		Content             string           `json:"content" doc:"Page content"`
		Metadata            BookPageMetadata `json:"metadata" doc:"Metadata of the book page"`
		AttachedImageFileID string           `json:"attachedImageFileId,omitempty" doc:"File ID for the attached image"`
	}
}

type CreateBookPageOutputBody struct {
	ID string `json:"id"`
}

type CreateBookPageOutput struct {
	Body CreateBookPageOutputBody
}

type BookPageOutput struct {
	ID                  string           `json:"id"`
	BookID              string           `json:"bookId"`
	PageNumber          int64            `json:"pageNumber"`
	Content             string           `json:"content"`
	Metadata            BookPageMetadata `json:"metadata"`
	AttachedImageFileID string           `json:"attachedImageFileId,omitempty"`
}

type GetBookPagesInput struct {
	BookID     string `query:"bookId" pattern:"^[a-fA-F0-9]{24}$" required:"true" doc:"Parent Book ID"`
	GetAll     bool   `query:"all" required:"true" default:"false" doc:"If true, returns all book pages ignoring pagination"`
	PageNumber int64  `query:"pageNumber" default:"1" minimum:"1" doc:"Page number"`
	PageSize   int64  `query:"pageSize" default:"50" minimum:"1" maximum:"500" doc:"Items per page"`
}

type GetBookPagesOutputBody struct {
	Data []BookPageOutput `json:"data"`
}

type GetBookPagesOutput struct {
	Body GetBookPagesOutputBody
}

type GetBookPagesRangeInput struct {
	BookID    string `query:"bookId" pattern:"^[a-fA-F0-9]{24}$" required:"true" doc:"Parent Book ID"`
	StartPage int64  `query:"startPage" required:"true" doc:"Start page number"`
	EndPage   int64  `query:"endPage" required:"true" doc:"End page number"`
}

type GetBookPagesOffsetInput struct {
	BookID     string `query:"bookId" pattern:"^[a-fA-F0-9]{24}$" required:"true" doc:"Book ID to query"`
	CenterPage int64  `query:"centerPage" required:"true" doc:"Center page to be offset from"`
	Offset     int64  `query:"offset" required:"true" default:"10" doc:"Number of frames before/after center"`
}

type GetBookPageByIDInput struct {
	ID string `path:"id" pattern:"^[a-fA-F0-9]{24}$" doc:"Entry ID"`
}

type GetBookPageByIDOutput struct {
	Body BookPageOutput
}
