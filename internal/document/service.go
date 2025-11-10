package internal

import (
	"context"
	"fmt"
)

type DocumentService struct {
	repo DocumentRepository
}

func (s *DocumentService) CreateNewDocument(ctx context.Context, data CreateDocumentDTO) (string, error) {
	docID, err := s.repo.CreateDocument(ctx, Document{
		Title:   data.Title,
		OwnerID: data.OwnerID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create a new document: %w", err)
	}
	return docID, nil
}

func (s *DocumentService) GetUserDocuments(ctx context.Context, ownerID string) ([]Document, error) {
	documents, err := s.repo.GetAllDocuments(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user documents: %w", err)
	}
	return documents, nil
}

func (s *DocumentService) GetOneDocument(ctx context.Context, data GetOneDocumentDTO) *Document {
	document := s.repo.FindDocument(ctx, data.OwnerID, data.DocumentID)
	return document
}

func NewDocumentService(documentsRepository DocumentRepository) (*DocumentService, error) {
	if documentsRepository == nil {
		return nil, fmt.Errorf("error creating the documents service: documentsRepository cannot be nil")
	}

	return &DocumentService{
		repo: documentsRepository,
	}, nil
}
