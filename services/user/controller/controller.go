package controller

import (
	"context"

	"github.com/charmbracelet/log"
	proto "github.com/vasapolrittideah/money-tracker-api/protogen"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"google.golang.org/grpc"
)

type userController struct {
	usecase domain.UserUsecase
	config  *config.Config
	proto.UnimplementedUserServiceServer
}

func NewUserController(grpc *grpc.Server, usecase domain.UserUsecase, config *config.Config) {
	handler := &userController{
		usecase: usecase,
		config:  config,
	}

	proto.RegisterUserServiceServer(grpc, handler)
}

func (c *userController) GetAllUsers(ctx context.Context, req *proto.Empty) (*proto.GetAllUsersResponse, error) {
	users, err := c.usecase.GetAllUsers()
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*proto.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &proto.GetAllUsersResponse{Users: protoUsers}, nil
}

func (c *userController) GetUserByID(
	ctx context.Context,
	req *proto.GetUserByIDRequest,
) (*proto.GetUserByIDResponse, error) {
	user, err := c.usecase.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GetUserByIDResponse{User: user.ToProto()}, nil
}

func (c *userController) GetUserByEmail(
	ctx context.Context,
	req *proto.GetUserByEmailRequest,
) (*proto.GetUserByEmailResponse, error) {
	user, err := c.usecase.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return &proto.GetUserByEmailResponse{User: user.ToProto()}, nil
}

func (c *userController) CreateUser(
	ctx context.Context,
	req *proto.CreateUserRequest,
) (*proto.CreateUserResponse, error) {
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

	return &proto.CreateUserResponse{User: createdUser.ToProto()}, nil
}

func (c *userController) UpdateUser(
	ctx context.Context,
	req *proto.UpdateUserRequest,
) (*proto.UpdateUserResponse, error) {
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

	return &proto.UpdateUserResponse{User: updatedUser.ToProto()}, nil
}

func (c *userController) DeleteUser(
	ctx context.Context,
	req *proto.DeleteUserRequest,
) (*proto.DeleteUserResponse, error) {
	deletedUser, err := c.usecase.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteUserResponse{User: deletedUser.ToProto()}, nil
}
