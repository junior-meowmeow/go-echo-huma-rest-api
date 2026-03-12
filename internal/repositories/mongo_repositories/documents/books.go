package documents

import (
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Name             string       `bson:"name"`
	Description      string       `bson:"description"`
	Metadata         BookMetadata `bson:"metadata"`
	CoverImageFileID string       `bson:"coverImageFileID"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

type BookMetadata struct {
	Author string `bson:"author"`
	ISBN   string `bson:"isbn"`
	Genre  string `bson:"genre"`
}

func NewBookDocument(entity *entities.Book) (BookDocument, error) {
	var bookDocument BookDocument
	var err error

	var oid bson.ObjectID
	if entity.ID != "" {
		oid, err = bson.ObjectIDFromHex(entity.ID)
		if err != nil {
			return bookDocument, fmt.Errorf("invalid book ID format: %w", err)
		}
	}

	bookDocument = BookDocument{
		ID:          oid,
		Name:        entity.Name,
		Description: entity.Description,
		Metadata: BookMetadata{
			Author: entity.Metadata.Author,
			ISBN:   entity.Metadata.ISBN,
			Genre:  entity.Metadata.Genre,
		},
		CoverImageFileID: entity.CoverImageFileID,
		CreatedAt:        entity.CreatedAt,
		ModifiedAt:       entity.ModifiedAt,
	}

	return bookDocument, nil
}

func (document *BookDocument) ToEntity() entities.Book {
	return entities.Book{
		ID:          document.ID.Hex(),
		Name:        document.Name,
		Description: document.Description,
		Metadata: entities.BookMetadata{
			Author: document.Metadata.Author,
			ISBN:   document.Metadata.ISBN,
			Genre:  document.Metadata.Genre,
		},
		CoverImageFileID: document.CoverImageFileID,
		CreatedAt:        document.CreatedAt,
		ModifiedAt:       document.ModifiedAt,
	}
}
