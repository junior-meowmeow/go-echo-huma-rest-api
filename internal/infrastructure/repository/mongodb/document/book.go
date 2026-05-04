package document

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type BookDocument struct {
	ID uuid.UUID `bson:"_id"`

	Name             string       `bson:"name"`
	Description      string       `bson:"description"`
	Metadata         BookMetadata `bson:"metadata"`
	CoverImageFileID string       `bson:"coverImageFileId"`

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

	bookUUID, err := StringToUUID(book.ID)
	if err != nil {
		return bookDocument, fmt.Errorf("invalid book ID format: %w", err)
	}

	bookDocument = BookDocument{
		ID:          bookUUID,
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
		ID:          doc.ID.String(),
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
