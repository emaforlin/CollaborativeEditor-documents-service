package repository

import (
	"context"

	"github.com/emaforlin/ce-document-service/models"
)

type DocumentRepository interface {
	CreateDocument(ctx context.Context, document models.Document) (string, error)
	GetAllDocuments(ctx context.Context, ownerID string) ([]models.Document, error)
	FindDocument(ctx context.Context, ownerID, documentID string) *models.Document
}
