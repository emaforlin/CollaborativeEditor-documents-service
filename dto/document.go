package dto

type CreateDocumentDTO struct {
	Title   string `json:"title" binding:"required"`
	OwnerID string
}

type GetOneDocumentDTO struct {
	DocumentID string `json:"document_id" binding:"required"`
	OwnerID    string
}
