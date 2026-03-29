package document

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type FileRecordDocument struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	FileName    string `bson:"fileName"`
	Size        int64  `bson:"size"`
	ContentType string `bson:"contentType"`
	S3Key       string `bson:"s3Key"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func NewFileRecordDocument(fileRecord *entity.FileRecord) (FileRecordDocument, error) {
	var fileRecordDocument FileRecordDocument
	var err error

	oid, err := StringToObjectID(fileRecord.ID)
	if err != nil {
		return fileRecordDocument, fmt.Errorf("invalid file record ID format: %w", err)
	}

	fileRecordDocument = FileRecordDocument{
		ID:          oid,
		FileName:    fileRecord.FileName,
		Size:        fileRecord.Size,
		ContentType: fileRecord.ContentType,
		S3Key:       fileRecord.S3Key,
		CreatedAt:   fileRecord.CreatedAt,
		UpdatedAt:   fileRecord.UpdatedAt,
	}

	return fileRecordDocument, nil
}

func (doc *FileRecordDocument) ToEntity() entity.FileRecord {
	return entity.FileRecord{
		ID:          doc.ID.Hex(),
		FileName:    doc.FileName,
		Size:        doc.Size,
		ContentType: doc.ContentType,
		S3Key:       doc.S3Key,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}
