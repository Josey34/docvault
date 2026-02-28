package entity

import "time"

type Document struct {
	ID          string
	FileName    string
	FileSize    int64
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	ContentType string
}
