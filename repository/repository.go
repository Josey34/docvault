package repository

import (
	"context"
	"docvault/entity"
)

type DocumentRepository interface {
	Save(ctx context.Context, doc *entity.Document) error
}
