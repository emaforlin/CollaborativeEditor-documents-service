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

// GetDocumentPermissions implements DocumentRepository
func (r *PostgresDocumentRepositoryImpl) GetDocumentPermissions(ctx context.Context, documentID string) []DocumentPermission {
	permissions, err := gorm.G[DocumentPermission](r.db).Where("document_id = ?", documentID).Find(ctx)
	if err != nil {
		return nil
	}
	return permissions
}

// RemoveDocumentPermission implements DocumentRepository
func (r *PostgresDocumentRepositoryImpl) RemoveDocumentPermission(ctx context.Context, userID, documentID string) error {
	if _, err := gorm.G[DocumentPermission](r.db).Where("document_id = ? AND user_id = ?", documentID, userID).Delete(ctx); err != nil {
		return fmt.Errorf("failed deleting document permission record: %w", err)
	}
	return nil
}

// CreateDocumentPermission implements DocumentRepository
func (r *PostgresDocumentRepositoryImpl) CreateDocumentPermission(ctx context.Context, permission DocumentPermission) error {
	if err := gorm.G[DocumentPermission](r.db).Create(ctx, &permission); err != nil {
		return fmt.Errorf("failed to create document permission: %w", err)
	}
	return nil
}

// UpdateDocument implements DocumentRepository
func (r *PostgresDocumentRepositoryImpl) UpdateDocument(ctx context.Context, document Document) error {
	updated, err := gorm.G[Document](r.db).Where("id = ?", document.ID).Updates(ctx, document)
	if updated < 1 {
		return fmt.Errorf("document update failed: no document matched")
	}

	if err != nil {
		return fmt.Errorf("document update failed: %w", err)
	}

	return nil
}

// CreateDocument implements DocumentRepository.
func (r *PostgresDocumentRepositoryImpl) CreateDocument(ctx context.Context, document Document) (string, error) {
	if err := gorm.G[Document](r.db).Create(ctx, &document); err != nil {
		return "", fmt.Errorf("failed to create document: %w", err)
	}
	return document.ID, nil
}

// FindDocument implements DocumentRepository.
func (r *PostgresDocumentRepositoryImpl) FindDocument(ctx context.Context, userID string, documentID string) *Document {
	var document Document

	// Query to find document if user is owner OR has permission
	err := r.db.WithContext(ctx).
		Select("DISTINCT documents.*").
		Table("documents").
		Joins("LEFT JOIN document_permissions ON documents.id = document_permissions.document_id").
		Where("documents.id = ? AND (documents.owner_id = ? OR document_permissions.user_id = ?)",
			documentID, userID, userID).
		First(&document).Error

	if err != nil {
		return nil
	}
	return &document
}

// GetAllDocuments implements DocumentRepository.
func (r *PostgresDocumentRepositoryImpl) GetUserDocuments(ctx context.Context, userID string, userIsOwner bool) ([]Document, error) {
	if userIsOwner {
		documents, err := gorm.G[Document](r.db).Where("owner_id = ?", userID).Find(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to find documents: %w", err)
		}
		return documents, nil
	}

	documents, err := gorm.G[Document](r.db).
		Raw("SELECT DISTINCT documents.* FROM documents LEFT JOIN document_permissions ON documents.id = document_permissions.document_id WHERE documents.owner_id = ? OR document_permissions.user_id = ?", userID, userID).
		Find(ctx)

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
