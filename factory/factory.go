package factory

import (
	"database/sql"
	"docvault/config"
	"docvault/database"
	"docvault/handler"
	"docvault/repository"
	"docvault/service"
	"docvault/usecase"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Factory struct {
	DB              *sql.DB
	DocumentHandler *handler.DocumentHandler
}

func New(cfg *config.Config) (*Factory, error) {
	db, err := database.NewSQLite(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to inititalize sqlite: %w", err)
	}

	if err := database.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("Error initializing minio client factory %w", err)
	}

	docRepo := repository.NewSQLiteDocumentRepository(db)

	storageService := service.NewMinIOStorage(minioClient, cfg.MinioBucketName)

	docUsecase := usecase.NewDocumentUsecase(docRepo, storageService)

	docHandler := handler.NewDocumentHandler(docUsecase)

	return &Factory{
		DB:              db,
		DocumentHandler: docHandler,
	}, nil
}
