package grpchandler

import (
	"context"

	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	user "github.com/vasapolrittideah/money-tracker-api/services/user/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/protoutil"
)

type userGRPCHandler struct {
	usecase user.UserUsecase
	config  *config.Config
	userpbv1.UnimplementedUserServiceServer
}

func NewUserGRPCHandler(usecase user.UserUsecase, config *config.Config) *userGRPCHandler {
	return &userGRPCHandler{
		usecase: usecase,
		config:  config,
	}
}

func (c *userGRPCHandler) GetAllUsers(
	ctx context.Context,
	req *userpbv1.GetAllUsersRequest,
) (*userpbv1.GetAllUsersResponse, error) {
	users, err := c.usecase.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*userpbv1.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, user.ToProto())
	}

	return &userpbv1.GetAllUsersResponse{Users: protoUsers}, nil
}

func (c *userGRPCHandler) GetUserByID(
	ctx context.Context,
	req *userpbv1.GetUserByIDRequest,
) (*userpbv1.GetUserByIDResponse, error) {
	user, err := c.usecase.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpbv1.GetUserByIDResponse{User: user.ToProto()}, nil
}

func (c *userGRPCHandler) GetUserByEmail(
	ctx context.Context,
	req *userpbv1.GetUserByEmailRequest,
) (*userpbv1.GetUserByEmailResponse, error) {
	user, err := c.usecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &userpbv1.GetUserByEmailResponse{User: user.ToProto()}, nil
}

func (c *userGRPCHandler) CreateUser(
	ctx context.Context,
	req *userpbv1.CreateUserRequest,
) (*userpbv1.CreateUserResponse, error) {
	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	}

	created, err := c.usecase.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userpbv1.CreateUserResponse{User: created.ToProto()}, nil
}

func (c *userGRPCHandler) UpdateUser(
	ctx context.Context,
	req *userpbv1.UpdateUserRequest,
) (*userpbv1.UpdateUserResponse, error) {
	updated, err := c.usecase.UpdateUser(ctx, req.Id, &user.UpdateUserRequest{
		FullName:   protoutil.UnwrapString(req.FullName),
		Email:      protoutil.UnwrapString(req.Email),
		Verified:   protoutil.UnwrapBool(req.Verified),
		Registered: protoutil.UnwrapBool(req.Registered),
		Password:   protoutil.UnwrapString(req.Password),
	})
	if err != nil {
		return nil, err
	}

	return &userpbv1.UpdateUserResponse{User: updated.ToProto()}, nil
}

func (c *userGRPCHandler) DeleteUser(
	ctx context.Context,
	req *userpbv1.DeleteUserRequest,
) (*userpbv1.DeleteUserResponse, error) {
	deleted, err := c.usecase.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpbv1.DeleteUserResponse{User: deleted.ToProto()}, nil
}
