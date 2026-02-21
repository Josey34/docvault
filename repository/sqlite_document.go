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

func (r *SQLiteDocumentRepository) FindById(ctx context.Context, id string) (*entity.Document, error) {
	findByIdQuery := `SELECT id, file_name, file_size, content_type, created_at, expires_at FROM documents WHERE id = ?`

	row := r.db.QueryRowContext(ctx, findByIdQuery, id)

	doc := &entity.Document{}
	err := row.Scan(&doc.ID, &doc.FileName, &doc.FileSize, &doc.ContentType, &doc.CreatedAt, &doc.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("error fetching document %w", err)
	}

	return doc, nil
}

func (r *SQLiteDocumentRepository) FindAll(ctx context.Context) ([]*entity.Document, error) {
	findAllQuery := `SELECT id, file_name, file_size, content_type, created_at, expires_at FROM documents ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, findAllQuery)
	if err != nil {
		return nil, fmt.Errorf("Error finding data from documents %w", err)
	}
	defer rows.Close()

	var documents []*entity.Document
	for rows.Next() {
		doc := &entity.Document{}
		err := rows.Scan(&doc.ID, &doc.FileName, &doc.FileSize, &doc.ContentType, &doc.CreatedAt, &doc.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning document %w", err)
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

func (r *SQLiteDocumentRepository) Delete(ctx context.Context, id string) error {
	deleteQuery := `DELETE FROM documents where id=?`

	_, err := r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("error deleting document %w", err)
	}

	return nil
}
