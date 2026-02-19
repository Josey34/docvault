package service

import (
	"context"
	"io"
)

type StorageService interface {
	Upload(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error
}
