package document

import (
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookPageDocument struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	BookID bson.ObjectID `bson:"book_id,omitempty"`

	PageNumber          int64            `bson:"pageNumber"`
	Content             string           `bson:"content"`
	Metadata            BookPageMetadata `bson:"metadata"`
	AttachedImageFileID string           `bson:"attachedImageFileID,omitempty"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

type BookPageMetadata struct {
	IsBookmarked bool   `bson:"isBookmarked"`
	Highlight    string `bson:"highlight"`
}

func NewBookPageDocument(entity *entity.BookPage) (BookPageDocument, error) {
	var bookPageDocument BookPageDocument

	bookOID, err := bson.ObjectIDFromHex(entity.BookID)
	if err != nil {
		return bookPageDocument, fmt.Errorf("invalid book ID format: %w", err)
	}

	var oid bson.ObjectID
	if entity.ID != "" {
		oid, err = bson.ObjectIDFromHex(entity.ID)
		if err != nil {
			return bookPageDocument, fmt.Errorf("invalid book page ID format: %w", err)
		}
	}

	bookPageDocument = BookPageDocument{
		ID:         oid,
		BookID:     bookOID,
		PageNumber: entity.PageNumber,
		Content:    entity.Content,
		Metadata: BookPageMetadata{
			IsBookmarked: entity.Metadata.IsBookmarked,
			Highlight:    entity.Metadata.Highlight,
		},
		AttachedImageFileID: entity.AttachedImageFileID,
		CreatedAt:           entity.CreatedAt,
		ModifiedAt:          entity.ModifiedAt,
	}

	return bookPageDocument, nil
}

func (document *BookPageDocument) ToEntity() entity.BookPage {
	return entity.BookPage{
		ID:         document.ID.Hex(),
		BookID:     document.BookID.Hex(),
		PageNumber: document.PageNumber,
		Content:    document.Content,
		Metadata: entity.BookPageMetadata{
			IsBookmarked: document.Metadata.IsBookmarked,
			Highlight:    document.Metadata.Highlight,
		},
		AttachedImageFileID: document.AttachedImageFileID,
		CreatedAt:           document.CreatedAt,
		ModifiedAt:          document.ModifiedAt,
	}
}
