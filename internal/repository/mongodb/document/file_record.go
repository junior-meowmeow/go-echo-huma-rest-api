package document

import (
	"fmt"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FileRecordDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	FileName    string `bson:"fileName"`
	Size        int64  `bson:"size"`
	ContentType string `bson:"contentType"`
	S3Key       string `bson:"s3Key"`

	CreatedAt  time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
}

func NewFileRecordDocument(entity *entity.FileRecord) (FileRecordDocument, error) {
	var fileRecordDocument FileRecordDocument
	var err error

	var oid bson.ObjectID
	if entity.ID != "" {
		oid, err = bson.ObjectIDFromHex(entity.ID)
		if err != nil {
			return fileRecordDocument, fmt.Errorf("invalid file record ID format: %w", err)
		}
	}

	fileRecordDocument = FileRecordDocument{
		ID:          oid,
		FileName:    entity.FileName,
		Size:        entity.Size,
		ContentType: entity.ContentType,
		S3Key:       entity.S3Key,
		CreatedAt:   entity.CreatedAt,
		ModifiedAt:  entity.ModifiedAt,
	}

	return fileRecordDocument, nil
}

func (document *FileRecordDocument) ToEntity() entity.FileRecord {
	return entity.FileRecord{
		ID:          document.ID.Hex(),
		FileName:    document.FileName,
		Size:        document.Size,
		ContentType: document.ContentType,
		S3Key:       document.S3Key,
		CreatedAt:   document.CreatedAt,
		ModifiedAt:  document.ModifiedAt,
	}
}
