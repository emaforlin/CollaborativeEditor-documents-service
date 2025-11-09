package core

import (
	"fmt"

	"github.com/emaforlin/ce-document-service/dto"
	"github.com/emaforlin/ce-document-service/models"
	"github.com/emaforlin/ce-document-service/repository"
)

type DocumentService struct {
	repo repository.DocumentRepository
}

func (s *DocumentService) CreateNewDocument(data dto.CreateDocumentDTO) (string, error) {
	docID, err := s.repo.CreateDocument(models.Document{
		Title:   data.Title,
		OwnerID: data.OwnerID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create a new document: %w", err)
	}
	return docID, nil
}

func (s *DocumentService) GetUserDocuments(ownerID string) ([]models.Document, error) {
	documents, err := s.repo.GetAllDocuments(ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user documents: %w", err)
	}
	return documents, nil
}

func (s *DocumentService) GetOneDocument(data dto.GetOneDocumentDTO) *models.Document {
	document := s.repo.FindDocument(data.OwnerID, data.DocumentID)
	return document
}

func NewDocumentService(documentsRepository repository.DocumentRepository) (*DocumentService, error) {
	if documentsRepository == nil {
		return nil, fmt.Errorf("error creating the documents service: documentsRepository cannot be nil")
	}

	return &DocumentService{
		repo: documentsRepository,
	}, nil
}
