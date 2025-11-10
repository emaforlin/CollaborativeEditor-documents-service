package internal

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const defaultTestDocumentJSON = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"This is a test document"}]}]}`

type MockRepository struct {
	documents []Document
}

// CreateDocument implements DocumentRepository.
func (m *MockRepository) CreateDocument(document Document) (documentID string, err error) {
	document.ID = uuid.NewString()
	document.Content = pgtype.JSONB{Bytes: []byte(defaultTestDocumentJSON), Status: pgtype.Present}
	m.documents = append(m.documents, document)
	return document.ID, nil
}

// GetAllDocuments implements DocumentRepository.
func (m *MockRepository) GetAllDocuments(ownerID string) ([]Document, error) {
	matchedDocs := make([]Document, 0)
	for _, doc := range m.documents {
		if doc.OwnerID == ownerID {
			matchedDocs = append(matchedDocs, doc)
		}
	}
	return matchedDocs, nil
}

// GetOneDocument implements DocumentRepository.
func (m *MockRepository) FindDocument(ownerID string, documentID string) *Document {
	for _, doc := range m.documents {
		if doc.OwnerID == ownerID && doc.ID == documentID {
			return &doc
		}
	}
	return nil
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		documents: []Document{
			{
				ID:        "mock-1",
				OwnerID:   "mock-fake-owner-1",
				Title:     "Mock 1 - Test document",
				Content:   pgtype.JSONB{Bytes: []byte(defaultTestDocumentJSON), Status: pgtype.Present},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "mock-2",
				OwnerID:   "mock-fake-owner-2",
				Title:     "Mock 2 - Collaboration notes",
				Content:   pgtype.JSONB{Bytes: []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Collaboration meeting notes here"}]}]}`), Status: pgtype.Present},
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-2 * time.Hour),
			},
			{
				ID:        "mock-3",
				OwnerID:   "mock-fake-owner-1",
				Title:     "Mock 3 - Meeting minutes",
				Content:   pgtype.JSONB{Bytes: []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Meeting minutes from last week"}]}]}`), Status: pgtype.Present},
				CreatedAt: time.Now().Add(-48 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				ID:        "mock-4",
				OwnerID:   "mock-fake-owner-3",
				Title:     "Mock 4 - Project plan",
				Content:   pgtype.JSONB{Bytes: []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Project planning document"}]}]}`), Status: pgtype.Present},
				CreatedAt: time.Now().Add(-72 * time.Hour),
				UpdatedAt: time.Now().Add(-6 * time.Hour),
			},
			{
				ID:        "mock-5",
				OwnerID:   "mock-fake-owner-2",
				Title:     "Mock 5 - Draft article",
				Content:   pgtype.JSONB{Bytes: []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Draft article content goes here"}]}]}`), Status: pgtype.Present},
				CreatedAt: time.Now().Add(-7 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-3 * 24 * time.Hour),
			},
			{
				ID:        "mock-6",
				OwnerID:   "mock-fake-owner-4",
				Title:     "Mock 6 - Research notes",
				Content:   pgtype.JSONB{Bytes: []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Research notes and findings"}]}]}`), Status: pgtype.Present},
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-29 * 24 * time.Hour),
			},
		},
	}
}
