package database

import (
	"time"

	"github.com/stretchr/testify/require"
)

type Adapter interface {
	CleanBooks(t require.TestingT)
	CountBooks(t require.TestingT) int64
	GetBookByID(t require.TestingT, id string) TestBookRecord
	SeedBooks(t require.TestingT, books []TestBookRecord)

	CleanBookPages(t require.TestingT)
	CountBookPages(t require.TestingT) int64
	GetBookPageByID(t require.TestingT, id string) TestBookPageRecord
	SeedBookPages(t require.TestingT, pages []TestBookPageRecord)

	CleanFiles(t require.TestingT)
	CountFiles(t require.TestingT) int64

	CleanReviews(t require.TestingT)
	CountReviews(t require.TestingT) int64
	GetAllReviews(t require.TestingT) []TestReviewRecord
	SeedReviews(t require.TestingT, reviews []TestReviewRecord)

	CleanUsers(t require.TestingT)
	CountUsers(t require.TestingT) int64
	GetUserByUsername(t require.TestingT, username string) TestUserRecord
}

type TestBookRecord struct {
	ID               string
	Name             string
	Description      string
	Author           string
	ISBN             string
	Genre            string
	CoverImageFileID string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type TestBookPageRecord struct {
	ID                  string
	BookID              string
	PageNumber          int64
	Content             string
	IsBookmarked        bool
	Highlight           string
	AttachedImageFileID string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type TestFileRecord struct {
	ID          string
	FileName    string
	Size        int64
	ContentType string
	S3Key       time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TestReviewRecord struct {
	ID        string
	Author    string
	Rating    int
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TestUserRecord struct {
	ID        string
	Username  string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
