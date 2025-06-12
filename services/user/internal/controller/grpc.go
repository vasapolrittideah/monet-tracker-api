package controller

import (
	"context"

	"github.com/charmbracelet/log"
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
)

type userGRPCController struct {
	usecase domain.UserUsecase
	config  *config.Config
	userpbv1.UnimplementedUserServiceServer
}

func NewUserGRPCController(usecase domain.UserUsecase, config *config.Config) *userGRPCController {
	return &userGRPCController{
		usecase: usecase,
		config:  config,
	}
}

func (c *userGRPCController) GetAllUsers(
	ctx context.Context,
	req *userpbv1.GetAllUsersRequest,
) (*userpbv1.GetAllUsersResponse, error) {
	users, err := c.usecase.GetAllUsers()
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*userpbv1.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &userpbv1.GetAllUsersResponse{Users: protoUsers}, nil
}

func (c *userGRPCController) GetUserByID(
	ctx context.Context,
	req *userpbv1.GetUserByIDRequest,
) (*userpbv1.GetUserByIDResponse, error) {
	user, err := c.usecase.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &userpbv1.GetUserByIDResponse{User: user.ToProto()}, nil
}

func (c *userGRPCController) GetUserByEmail(
	ctx context.Context,
	req *userpbv1.GetUserByEmailRequest,
) (*userpbv1.GetUserByEmailResponse, error) {
	user, err := c.usecase.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return &userpbv1.GetUserByEmailResponse{User: user.ToProto()}, nil
}

func (c *userGRPCController) CreateUser(
	ctx context.Context,
	req *userpbv1.CreateUserRequest,
) (*userpbv1.CreateUserResponse, error) {
	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := c.usecase.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &userpbv1.CreateUserResponse{User: createdUser.ToProto()}, nil
}

func (c *userGRPCController) UpdateUser(
	ctx context.Context,
	req *userpbv1.UpdateUserRequest,
) (*userpbv1.UpdateUserResponse, error) {
	user := &domain.User{
		FullName:     req.FullName,
		Email:        req.Email,
		Verified:     req.Verified,
		Password:     req.Password,
		RefreshToken: req.RefreshToken,
	}

	updatedUser, err := c.usecase.UpdateUser(req.Id, user)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &userpbv1.UpdateUserResponse{User: updatedUser.ToProto()}, nil
}

func (c *userGRPCController) DeleteUser(
	ctx context.Context,
	req *userpbv1.DeleteUserRequest,
) (*userpbv1.DeleteUserResponse, error) {
	deletedUser, err := c.usecase.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &userpbv1.DeleteUserResponse{User: deletedUser.ToProto()}, nil
}
