package domain

import (
	"time"

	"github.com/google/uuid"
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID         uuid.UUID `json:"id"         gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FullName   string    `json:"full_name"  gorm:"not null;type:varchar(100)"                      example:"John Doe"             extensions:"x-order=2"`
	Email      string    `json:"email"      gorm:"not null;uniqueIndex"                            example:"john@example.com"     extensions:"x-order=3"`
	Verified   bool      `json:"verified"   gorm:"not null;default:false"                          example:"true"                 extensions:"x-order=4"`
	Registered bool      `json:"registered" gorm:"not null;default:false"                          example:"true"                 extensions:"x-order=5"`
	Password   string    `json:"-"          gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"                                  example:"2022-01-01T00:00:00Z" extensions:"x-order=6"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"                                  example:"2022-01-01T00:00:00Z" extensions:"x-order=7"`
	Sessions   []Session `json:"-"          gorm:"foreignKey:UserID"`
}

func (u *User) ToProto() *userpbv1.User {
	return &userpbv1.User{
		Id:         u.ID.String(),
		FullName:   u.FullName,
		Email:      u.Email,
		Verified:   u.Verified,
		Registered: u.Registered,
		Password:   u.Password,
		CreatedAt:  timestamppb.New(u.CreatedAt),
		UpdatedAt:  timestamppb.New(u.UpdatedAt),
	}
}

func NewUserFromProto(user *userpbv1.User) *User {
	return &User{
		ID:         uuid.MustParse(user.Id),
		FullName:   user.FullName,
		Email:      user.Email,
		Verified:   user.Verified,
		Registered: user.Registered,
		Password:   user.Password,
		CreatedAt:  user.CreatedAt.AsTime(),
		UpdatedAt:  user.UpdatedAt.AsTime(),
	}
}
