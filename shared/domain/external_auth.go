package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExternalAuth struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" `
	UserID     uuid.UUID `gorm:"not null"`
	Provider   string    `gorm:"not null"`
	ProviderID string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
