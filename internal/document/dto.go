package internal

type AddCollaboratorDTO struct {
	DocumentID     string
	OwnerID        string
	CollaboratorID string `json:"user_id"`
	Role           Role   `json:"role"`
}

type CreateDocumentDTO struct {
	Title   string `json:"title" binding:"required"`
	OwnerID string
}

type GetOneDocumentDTO struct {
	DocumentID string `json:"document_id" binding:"required"`
	OwnerID    string
}
