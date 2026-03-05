package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookMetadata struct {
	Author string `bson:"author"`
	ISBN   string `bson:"isbn"`
	Genre  string `bson:"genre"`
}

type Book struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Name             string       `bson:"name"`
	Description      string       `bson:"description"`
	Metadata         BookMetadata `bson:"metadata"`
	CoverImageFileID string       `bson:"coverImageFileID"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}
