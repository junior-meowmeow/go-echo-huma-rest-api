package documents

import (
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FileMetadataDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Filename    string `bson:"filename"`
	Size        int64  `bson:"size"`
	ContentType string `bson:"contentType"`
	S3Key       string `bson:"s3Key"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func NewFileMetadataDocument(entity *entities.FileMetadata) (FileMetadataDocument, error) {
	var fileMetadataDocument FileMetadataDocument
	var err error

	var oid bson.ObjectID
	if entity.ID != "" {
		oid, err = bson.ObjectIDFromHex(entity.ID)
		if err != nil {
			return fileMetadataDocument, fmt.Errorf("invalid file metadata ID format: %w", err)
		}
	}

	fileMetadataDocument = FileMetadataDocument{
		ID:          oid,
		Filename:    entity.Filename,
		Size:        entity.Size,
		ContentType: entity.ContentType,
		S3Key:       entity.S3Key,
		CreatedAt:   entity.CreatedAt,
		ModifiedAt:  entity.ModifiedAt,
	}

	return fileMetadataDocument, nil
}

func (document *FileMetadataDocument) ToEntity() entities.FileMetadata {
	return entities.FileMetadata{
		ID:          document.ID.Hex(),
		Filename:    document.Filename,
		Size:        document.Size,
		ContentType: document.ContentType,
		S3Key:       document.S3Key,
		CreatedAt:   document.CreatedAt,
		ModifiedAt:  document.ModifiedAt,
	}
}
