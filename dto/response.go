package dto

import (
	"docvault/entity"
	"time"
)

type DocumentResponse struct {
	ID          string     `json:"id"`
	FileName    string     `json:"file_name"`
	FileSize    int64      `json:"file_size"`
	ContentType string     `json:"content_type"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

func FromEntity(doc *entity.Document) *DocumentResponse {
	return &DocumentResponse{
		ID:          doc.ID,
		FileName:    doc.FileName,
		FileSize:    doc.FileSize,
		ContentType: doc.ContentType,
		CreatedAt:   doc.CreatedAt,
		ExpiresAt:   doc.ExpiresAt,
	}
}
