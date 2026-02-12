package internal

import "time"

type UpdateDocumentDTO struct {
	DocumentID string
	Title      string `json:"title"`
}

type RemoveCollaboratorDTO struct {
	DocumentID string
	UserID     string `json:"user_id"`
}

type AddCollaboratorDTO struct {
	DocumentID string
	OwnerID    string
	UserID     string `json:"user_id"`
	Role       Role   `json:"role"`
}

type CreateDocumentDTO struct {
	Title   string `json:"title" binding:"required"`
	OwnerID string
}

type GetOneDocumentDTO struct {
	DocumentID string `json:"document_id" binding:"required"`
	OwnerID    string
}

type DocumentResponse struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentDetailResponse struct {
	ID        string      `json:"id"`
	OwnerID   string      `json:"owner_id"`
	Title     string      `json:"title"`
	Content   interface{} `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CollaboratorResponse struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
}

// Converter functions to transform models to DTOs

func ToDocumentResponse(doc *Document) DocumentResponse {
	return DocumentResponse{
		ID:        doc.ID,
		OwnerID:   doc.OwnerID,
		Title:     doc.Title,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

func ToDocumentDetailResponse(doc *Document) DocumentDetailResponse {
	var content interface{}
	if doc.Content != nil {
		// Convert pgtype.JSONB to interface{}
		if err := doc.Content.AssignTo(&content); err != nil {
			content = nil
		}
	}

	return DocumentDetailResponse{
		ID:        doc.ID,
		OwnerID:   doc.OwnerID,
		Title:     doc.Title,
		Content:   content,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}

func ToCollaboratorResponse(perm *DocumentPermission) CollaboratorResponse {
	return CollaboratorResponse{
		UserID: perm.UserID,
		Role:   perm.Role,
	}
}

// Generic function to convert a slice of models to a slice of responses
func ToResponseList[T any, R any](items []T, converter func(*T) R) []R {
	responses := make([]R, len(items))
	for i := range items {
		responses[i] = converter(&items[i])
	}
	return responses
}

// Convenience wrapper functions using the generic converter
func ToDocumentResponseList(docs []Document) []DocumentResponse {
	return ToResponseList(docs, ToDocumentResponse)
}

func ToCollaboratorResponseList(perms []DocumentPermission) []CollaboratorResponse {
	return ToResponseList(perms, ToCollaboratorResponse)
}
