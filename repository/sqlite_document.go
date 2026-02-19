package repository

import (
	"context"
	"database/sql"
	"docvault/entity"
	"fmt"
)

type SQLiteDocumentRepository struct {
	db *sql.DB
}

func NewSQLiteDocumentRepository(db *sql.DB) DocumentRepository {
	return &SQLiteDocumentRepository{db: db}
}

func (r *SQLiteDocumentRepository) Save(ctx context.Context, doc *entity.Document) error {
	insertQuery := `INSERT INTO documents (id, file_name, file_size, content_type, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, insertQuery, doc.ID, doc.FileName, doc.FileSize, doc.ContentType, doc.CreatedAt, doc.ExpiresAt)
	if err != nil {
		return fmt.Errorf("Error inserting data to database documents %w", err)
	}

	return nil
}
