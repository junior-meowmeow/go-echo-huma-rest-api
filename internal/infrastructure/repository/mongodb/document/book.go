package document

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type BookDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Name             string       `bson:"name"`
	Description      string       `bson:"description"`
	Metadata         BookMetadata `bson:"metadata"`
	CoverImageFileID string       `bson:"coverImageFileID"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type BookMetadata struct {
	Author string `bson:"author"`
	ISBN   string `bson:"isbn"`
	Genre  string `bson:"genre"`
}

func NewBookDocument(book *entity.Book) (BookDocument, error) {
	var bookDocument BookDocument
	var err error

	oid, err := StringToObjectID(book.ID)
	if err != nil {
		return bookDocument, fmt.Errorf("invalid book ID format: %w", err)
	}

	bookDocument = BookDocument{
		ID:          oid,
		Name:        book.Name,
		Description: book.Description,
		Metadata: BookMetadata{
			Author: book.Metadata.Author,
			ISBN:   book.Metadata.ISBN,
			Genre:  book.Metadata.Genre,
		},
		CoverImageFileID: book.CoverImageFileID,
		CreatedAt:        book.CreatedAt,
		UpdatedAt:        book.UpdatedAt,
	}

	return bookDocument, nil
}

func (doc *BookDocument) ToEntity() entity.Book {
	return entity.Book{
		ID:          doc.ID.Hex(),
		Name:        doc.Name,
		Description: doc.Description,
		Metadata: entity.BookMetadata{
			Author: doc.Metadata.Author,
			ISBN:   doc.Metadata.ISBN,
			Genre:  doc.Metadata.Genre,
		},
		CoverImageFileID: doc.CoverImageFileID,
		CreatedAt:        doc.CreatedAt,
		UpdatedAt:        doc.UpdatedAt,
	}
}
