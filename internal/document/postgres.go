package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/emaforlin/ce-document-service/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDocumentRepositoryImpl struct {
	db *gorm.DB
}

// CreateDocument implements DocumentRepository.
func (r *PostgresDocumentRepositoryImpl) CreateDocument(ctx context.Context, document Document) (string, error) {
	if err := gorm.G[Document](r.db).Create(ctx, &document); err != nil {
		return "", fmt.Errorf("failed to create document: %w", err)
	}
	return document.ID, nil
}

// FindDocument implements DocumentRepository.
func (r *PostgresDocumentRepositoryImpl) FindDocument(ctx context.Context, ownerID string, documentID string) *Document {
	document, err := gorm.G[Document](r.db).Where("id = ?, owner_id = ?", documentID, ownerID).First(ctx)
	if err != nil {
		return nil
	}
	return &document
}

// GetAllDocuments implements DocumentRepository.
func (r *PostgresDocumentRepositoryImpl) GetAllDocuments(ctx context.Context, ownerID string) ([]Document, error) {
	documents, err := gorm.G[Document](r.db).Where("owner_id = ?", ownerID).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	return documents, nil
}

func (r *PostgresDocumentRepositoryImpl) GetDB() *gorm.DB {
	return r.db
}

func NewPostgresRepository(cfg config.DatabaseConfig) *PostgresDocumentRepositoryImpl {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", cfg.Host, cfg.User, cfg.Pass, cfg.Name, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed database connection: ", err)
	}

	return &PostgresDocumentRepositoryImpl{
		db: db,
	}
}
