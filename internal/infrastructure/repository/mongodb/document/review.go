//nolint:dupl // Documents are intended to follow a similar pattern.
package document

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type ReviewDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Author  string `bson:"author"`
	Rating  int    `bson:"rating"`
	Message string `bson:"message"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func NewReviewDocument(review *entity.Review) (ReviewDocument, error) {
	var reviewDocument ReviewDocument
	var err error

	oid, err := StringToObjectID(review.ID)
	if err != nil {
		return reviewDocument, fmt.Errorf("invalid review ID format: %w", err)
	}

	reviewDocument = ReviewDocument{
		ID:        oid,
		Author:    review.Author,
		Rating:    review.Rating,
		Message:   review.Message,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}

	return reviewDocument, nil
}

func (doc *ReviewDocument) ToEntity() entity.Review {
	return entity.Review{
		ID:        doc.ID.Hex(),
		Author:    doc.Author,
		Rating:    doc.Rating,
		Message:   doc.Message,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
