package repository

import (
	"github.com/emaforlin/ce-document-service/models"
)

type DocumentRepository interface {
	CreateDocument(document models.Document) (documentID string, err error)
	GetAllDocuments(ownerID string) ([]models.Document, error)
	GetOneDocument(ownerID, documentID string) *models.Document
}
