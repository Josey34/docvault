package factory

import (
	"context"
	"database/sql"
	"docvault/config"
	"docvault/database"
	"docvault/handler"
	"docvault/repository"
	"docvault/service"
	"docvault/usecase"
	"docvault/worker"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Factory struct {
	DB                 *sql.DB
	DocumentHandler    *handler.DocumentHandler
	NotificationWorker *worker.NotificationWorker
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

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)

	_, err = sqsClient.CreateQueue(context.Background(), &sqs.CreateQueueInput{
		QueueName: aws.String("docvault-events"),
	})
	if err != nil {
		fmt.Printf("Queue creation warning: %v (might already exist)\n", err)
	}

	queueService := service.NewSQSQueue(sqsClient, cfg.SqsQueueUrl)

	docRepo := repository.NewSQLiteDocumentRepository(db)

	storageService := service.NewMinIOStorage(minioClient, cfg.MinioBucketName)

	docUsecase := usecase.NewDocumentUsecase(docRepo, storageService, queueService)

	docHandler := handler.NewDocumentHandler(docUsecase)

	notificationWorker := worker.NewNotificationWorker(queueService)

	return &Factory{
		DB:                 db,
		DocumentHandler:    docHandler,
		NotificationWorker: notificationWorker,
	}, nil
}
