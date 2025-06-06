package domain

import (
	proto "github.com/vasapolrittideah/money-tracker-api/protogen"
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

func (u *User) ToProto() *proto.User {
	return &proto.User{
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

// validate input body

type CreateUserInput struct {
	FullName           string `validate:"required"`
	Email              string `validate:"required,email"`
	Verified           bool
	HashedPassword     string `validate:"required"`
	HashedRefreshToken string
}

func (i *CreateUserInput) Populate(req *proto.CreateUserRequest) {
	i.FullName = req.FullName
	i.Email = req.Email
	i.HashedPassword = req.HashedPassword
}

type UpdateUserInput struct {
	ID                 uint64 `validate:"required"`
	FullName           string
	Email              string `validate:"email"`
	Verified           bool
	HashedPassword     string
	HashedRefreshToken string
}

func (i *UpdateUserInput) Populate(req *proto.UpdateUserRequest) {
	i.ID = req.Id
	i.FullName = req.FullName
	i.Email = req.Email
	i.HashedPassword = req.HashedPassword
	i.HashedRefreshToken = req.HashedRefreshToken
}
