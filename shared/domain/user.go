package domain

import (
	userv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName           string `gorm:"not null;type:varchar(100)"`
	Email              string `gorm:"not null;uniqueIndex"`
	Verified           bool   `gorm:"not null;default:false"`
	HashedPassword     string
	HashedRefreshToken string
}

func (u *User) ToProto() *userv1.User {
	return &userv1.User{
		Id:       uint64(u.ID),
		FullName: u.FullName,
		Email:    u.Email,
		Verified: u.Verified,
	}
}

type UserRepository interface {
	GetAllUsers() ([]*User, error)
	GetUserByID(id uint64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(id uint64, user *User) (*User, error)
	DeleteUser(id uint64) (*User, error)
}

type UserUsecase interface {
	GetAllUsers() ([]*User, error)
	GetUserByID(id uint64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(id uint64, user *User) (*User, error)
	DeleteUser(id uint64) (*User, error)
}
