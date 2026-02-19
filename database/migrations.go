package database

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	if err := CreateDocumentsTable(db); err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	if err := CreateUsersTable(db); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if err := CreateProcessingResultsTable(db); err != nil {
		return fmt.Errorf("failed to create processing results table: %w", err)
	}

	if err := CreateDocumentChunksTable(db); err != nil {
		return fmt.Errorf("failed to create document chunks table: %w", err)
	}

	return nil
}

func CreateDocumentsTable(db *sql.DB) error {
	createDocumentsQuery := ` CREATE TABLE IF NOT EXISTS documents (
        id TEXT PRIMARY KEY,
        file_name TEXT,
				file_size INTEGER,
				content_type TEXT,
				created_at DATETIME NOT NULL,
				expires_at DATETIME
    );
	`

	_, err := db.Exec(createDocumentsQuery)
	if err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	fmt.Println("Table 'documents' created successfully")
	return nil
}

func CreateUsersTable(db *sql.DB) error {
	createUsersQuery := ` CREATE TABLE IF NOT EXISTS users (
                    id TEXT PRIMARY KEY,
            username TEXT NOT NULL UNIQUE,
            email TEXT NOT NULL UNIQUE,
            created_at DATETIME NOT NULL
    );
	`
	_, err := db.Exec(createUsersQuery)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	fmt.Println("Table 'users' created successfully")

	return nil
}

func CreateProcessingResultsTable(db *sql.DB) error {
	createProcessingResultsQuery := ` CREATE TABLE IF NOT EXISTS processing_results (
            id TEXT PRIMARY KEY,
            document_id TEXT NOT NULL,
            text_content TEXT NOT NULL,
            created_at DATETIME NOT NULL
    );
	`

	_, err := db.Exec(createProcessingResultsQuery)
	if err != nil {
		return fmt.Errorf("failed to create processing_results table: %w", err)
	}

	fmt.Println("Table 'processing_results' created successfully")
	return nil
}

func CreateDocumentChunksTable(db *sql.DB) error {
	createDocumentChunksQuery := ` CREATE TABLE IF NOT EXISTS document_chunks (
            id TEXT PRIMARY KEY,
            document_id TEXT NOT NULL,
            chunk_number INTEGER NOT NULL,
						chunk_text TEXT NOT NULL,
            created_at DATETIME NOT NULL
    );
	`

	_, err := db.Exec(createDocumentChunksQuery)
	if err != nil {
		return fmt.Errorf("failed to create document_chunks table: %w", err)
	}

	fmt.Println("Table 'document_chunks' created successfully")
	return nil
}
