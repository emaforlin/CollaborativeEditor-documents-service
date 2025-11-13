package internal

import (
	"context"
)

type DocumentRepository interface {
	GetDocumentPermissions(ctx context.Context, documentID string) []DocumentPermission
	RemoveDocumentPermission(ctx context.Context, userID, documentID string) error
	CreateDocumentPermission(ctx context.Context, permission DocumentPermission) error

	DeleteDocument(ctx context.Context, documentID string) error
	UpdateDocument(ctx context.Context, document Document) error
	CreateDocument(ctx context.Context, document Document) (string, error)
	GetUserDocuments(ctx context.Context, userID string, userIsOwner bool) ([]Document, error)
	FindDocument(ctx context.Context, userID, documentID string) *Document

	GetDocumentWithPermission(ctx context.Context, userID, documentID string) (*Document, string)
}
