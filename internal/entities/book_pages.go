package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookPageMetadata struct {
	IsBookmarked bool   `bson:"isBookmarked"`
	Highlight    string `bson:"highlight"`
}

type BookPage struct {
	ID     bson.ObjectID `bson:"_id,omitempty"`
	BookID bson.ObjectID `bson:"book_id,omitempty"`

	PageNumber          int64            `bson:"pageNumber"`
	Content             string           `bson:"content"`
	Metadata            BookPageMetadata `bson:"metadata"`
	AttachedImageFileID string           `bson:"attachedImageFileID,omitempty"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}
