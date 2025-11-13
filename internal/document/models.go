package internal

import (
	"time"

	"github.com/jackc/pgtype"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleEditor Role = "editor"
	RoleViewer Role = "viewer"
)

type Document struct {
	ID            string               `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	OwnerID       string               `gorm:"type:uuid"`
	Title         string               `gorm:"size:255"`
	Content       *pgtype.JSONB        `gorm:"type:jsonb"`
	Collaborators []DocumentPermission `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DocumentPermission struct {
	ID         string `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	DocumentID string `gorm:"type:uuid;not null;uniqueIndex:idx_document_user_permission"`
	UserID     string `gorm:"type:uuid;not null;uniqueIndex:idx_document_user_permission"`
	Role       Role   `gorm:"type:varchar(10);not null"`
}
