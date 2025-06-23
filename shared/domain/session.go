package domain

import "time"

type Session struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	UserID    uint64 `gorm:"not null"`
	Token     string `gorm:"not null;unique"`
	UserAgent string
	IPAddress string
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
