package repository

import (
	"context"
	"docvault/entity"
	"time"
)

type DocumentRepository interface {
	Save(ctx context.Context, doc *entity.Document) error
	FindById(ctx context.Context, id string) (*entity.Document, error)
	FindAll(ctx context.Context) ([]*entity.Document, error)
	Delete(ctx context.Context, id string) error
	FindExpired(ctx context.Context, now time.Time) ([]*entity.Document, error)
}
