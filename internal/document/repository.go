package internal

import (
	"context"
)

type DocumentRepository interface {
	GetDocumentPermissions(ctx context.Context, documentID string) []DocumentPermission
	RemoveDocumentPermission(ctx context.Context, userID, documentID string) error
	CreateDocumentPermission(ctx context.Context, permission DocumentPermission) error

	UpdateDocument(ctx context.Context, document Document) error
	CreateDocument(ctx context.Context, document Document) (string, error)
	GetAllDocuments(ctx context.Context, ownerID string) ([]Document, error)
	FindDocument(ctx context.Context, ownerID, documentID string) *Document
}
