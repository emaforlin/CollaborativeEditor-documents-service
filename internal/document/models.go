package internal

import (
	"time"

	"github.com/jackc/pgtype"
)

type Document struct {
	ID        string       `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	OwnerID   string       `gorm:"type:uuid"`
	Title     string       `gorm:"size:255"`
	Content   pgtype.JSONB `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
