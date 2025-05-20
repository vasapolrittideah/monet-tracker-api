package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                 uuid.UUID `json:"id"              gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FullName           string    `json:"name"            gorm:"type:varchar(255)"`
	Email              string    `json:"email"           gorm:"uniqueIndex:not null;type:varchar(255)"`
	Verified           bool      `json:"verified"        gorm:"not null;default:false"`
	CreatedAt          time.Time `json:"created_at"      gorm:"autoCreateTime;not null"`
	UpdatedAt          time.Time `json:"updated_at"      gorm:"autoUpdateTime;not null"`
	LastSignInAt       time.Time `json:"last_sign_in_at"`
	HashedPassword     string    `json:"-"               gorm:"not null;type:varchar(255)"`
	HashedRefreshToken string    `json:"-"`
}
