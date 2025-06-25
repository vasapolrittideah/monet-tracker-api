package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" `
	UserID    uuid.UUID `gorm:"not null"`
	Token     string    `gorm:"not null;unique"`
	UserAgent string
	IPAddress string
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
