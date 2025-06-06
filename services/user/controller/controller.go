package controller

import (
	"context"

	"github.com/charmbracelet/log"
	userv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"google.golang.org/grpc"
)

type userController struct {
	usecase domain.UserUsecase
	config  *config.Config
	userv1.UnimplementedUserServiceServer
}

func NewUserController(grpc *grpc.Server, usecase domain.UserUsecase, config *config.Config) {
	handler := &userController{
		usecase: usecase,
		config:  config,
	}

	userv1.RegisterUserServiceServer(grpc, handler)
}

func (c *userController) GetAllUsers(
	ctx context.Context,
	req *userv1.GetAllUsersRequest,
) (*userv1.GetAllUsersResponse, error) {
	users, err := c.usecase.GetAllUsers()
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*userv1.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &userv1.GetAllUsersResponse{Users: protoUsers}, nil
}

func (c *userController) GetUserByID(
	ctx context.Context,
	req *userv1.GetUserByIDRequest,
) (*userv1.GetUserByIDResponse, error) {
	user, err := c.usecase.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserByIDResponse{User: user.ToProto()}, nil
}

func (c *userController) GetUserByEmail(
	ctx context.Context,
	req *userv1.GetUserByEmailRequest,
) (*userv1.GetUserByEmailResponse, error) {
	user, err := c.usecase.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserByEmailResponse{User: user.ToProto()}, nil
}

func (c *userController) CreateUser(
	ctx context.Context,
	req *userv1.CreateUserRequest,
) (*userv1.CreateUserResponse, error) {
	user := &domain.User{
		FullName:           req.FullName,
		Email:              req.Email,
		Verified:           req.Verified,
		HashedPassword:     req.HashedPassword,
		HashedRefreshToken: req.HashedRefreshToken,
	}

	createdUser, err := c.usecase.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &userv1.CreateUserResponse{User: createdUser.ToProto()}, nil
}

func (c *userController) UpdateUser(
	ctx context.Context,
	req *userv1.UpdateUserRequest,
) (*userv1.UpdateUserResponse, error) {
	user := &domain.User{
		FullName:           req.FullName,
		Email:              req.Email,
		Verified:           req.Verified,
		HashedPassword:     req.HashedPassword,
		HashedRefreshToken: req.HashedRefreshToken,
	}

	updatedUser, err := c.usecase.UpdateUser(req.Id, user)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &userv1.UpdateUserResponse{User: updatedUser.ToProto()}, nil
}

func (c *userController) DeleteUser(
	ctx context.Context,
	req *userv1.DeleteUserRequest,
) (*userv1.DeleteUserResponse, error) {
	deletedUser, err := c.usecase.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &userv1.DeleteUserResponse{User: deletedUser.ToProto()}, nil
}
