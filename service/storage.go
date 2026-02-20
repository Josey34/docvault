package service

import (
	"context"
	"io"
)

type StorageService interface {
	Upload(ctx context.Context, filename string, fileSize int64, contentType string, file io.Reader) error
	Download(ctx context.Context, filename string) (io.ReadCloser, error)
}
