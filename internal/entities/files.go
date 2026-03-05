package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FileMetadata struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Filename    string `bson:"filename"`
	Size        int64  `bson:"size"`
	ContentType string `bson:"contentType"`
	S3Key       string `bson:"s3Key"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}
