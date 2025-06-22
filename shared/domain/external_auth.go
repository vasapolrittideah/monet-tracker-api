package domain

import "time"

type ExternalAuth struct {
	ID         uint64 `json:"id"         gorm:"primaryKey;autoIncrement"`
	Provider   string `                  gorm:"not null"`
	ProviderID string `                  gorm:"not null"`
	UserID     uint64
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
