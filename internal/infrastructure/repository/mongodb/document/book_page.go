package document

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type BookPageDocument struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	BookID bson.ObjectID `bson:"book_id,omitempty"`

	PageNumber          int64            `bson:"pageNumber"`
	Content             string           `bson:"content"`
	Metadata            BookPageMetadata `bson:"metadata"`
	AttachedImageFileID string           `bson:"attachedImageFileId,omitempty"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type BookPageMetadata struct {
	IsBookmarked bool   `bson:"isBookmarked"`
	Highlight    string `bson:"highlight"`
}

func NewBookPageDocument(bookPage *entity.BookPage) (BookPageDocument, error) {
	var bookPageDocument BookPageDocument

	bookOID, err := StringToObjectID(bookPage.BookID)
	if err != nil {
		return bookPageDocument, fmt.Errorf("invalid book ID format: %w", err)
	}

	oid, err := StringToObjectID(bookPage.ID)
	if err != nil {
		return bookPageDocument, fmt.Errorf("invalid book page ID format: %w", err)
	}

	bookPageDocument = BookPageDocument{
		ID:         oid,
		BookID:     bookOID,
		PageNumber: bookPage.PageNumber,
		Content:    bookPage.Content,
		Metadata: BookPageMetadata{
			IsBookmarked: bookPage.Metadata.IsBookmarked,
			Highlight:    bookPage.Metadata.Highlight,
		},
		AttachedImageFileID: bookPage.AttachedImageFileID,
		CreatedAt:           bookPage.CreatedAt,
		UpdatedAt:           bookPage.UpdatedAt,
	}

	return bookPageDocument, nil
}

func (doc *BookPageDocument) ToEntity() entity.BookPage {
	return entity.BookPage{
		ID:         doc.ID.Hex(),
		BookID:     doc.BookID.Hex(),
		PageNumber: doc.PageNumber,
		Content:    doc.Content,
		Metadata: entity.BookPageMetadata{
			IsBookmarked: doc.Metadata.IsBookmarked,
			Highlight:    doc.Metadata.Highlight,
		},
		AttachedImageFileID: doc.AttachedImageFileID,
		CreatedAt:           doc.CreatedAt,
		UpdatedAt:           doc.UpdatedAt,
	}
}
