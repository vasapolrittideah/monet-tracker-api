//nolint:lll

package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                 uuid.UUID       `json:"id"              gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FullName           string          `json:"name"            gorm:"not null;type:varchar(100)"`
	Email              string          `json:"email"           gorm:"not null;uniqueIndex"`
	Verified           bool            `json:"verified"        gorm:"not null;default:false"`
	CreatedAt          time.Time       `json:"created_at"      gorm:"autoCreateTime"`
	UpdatedAt          time.Time       `json:"updated_at"      gorm:"autoUpdateTime"`
	LastSignInAt       *time.Time      `json:"last_sign_in_at"`
	HashedPassword     string          `json:"-"               gorm:"not null"`
	HashedRefreshToken string          `json:"-"`
	SocialAccounts     []SocialAccount `json:"-"               gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SocialAccount struct {
	Id         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Provider   string    `gorm:"not null"`
	ProviderId string    `gorm:"not null"`
	UserId     uuid.UUID `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

type Jwt struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
