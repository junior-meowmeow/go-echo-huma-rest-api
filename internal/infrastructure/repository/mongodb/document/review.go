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

func NewReviewDocument(entity *entity.Review) (ReviewDocument, error) {
	var reviewDocument ReviewDocument
	var err error

	var oid bson.ObjectID
	if entity.ID != "" {
		oid, err = bson.ObjectIDFromHex(entity.ID)
		if err != nil {
			return reviewDocument, fmt.Errorf("invalid review ID format: %w", err)
		}
	}

	reviewDocument = ReviewDocument{
		ID:        oid,
		Author:    entity.Author,
		Rating:    entity.Rating,
		Message:   entity.Message,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return reviewDocument, nil
}

func (document *ReviewDocument) ToEntity() entity.Review {
	return entity.Review{
		ID:        document.ID.Hex(),
		Author:    document.Author,
		Rating:    document.Rating,
		Message:   document.Message,
		CreatedAt: document.CreatedAt,
		UpdatedAt: document.UpdatedAt,
	}
}
