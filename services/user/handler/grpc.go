package handler

import (
	"context"

	"github.com/google/uuid"
	userpb "github.com/vasapolrittideah/money-tracker-api/generated/protobuf/user"
	"github.com/vasapolrittideah/money-tracker-api/services/user/service"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/mapper"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type UserGrpcHandler interface {
	GetAllUsers(c context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error)
	GetUserById(c context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserByIdResponse, error)
	GetUserByEmail(c context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error)
	CreateUser(c context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
	UpdateUser(c context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error)
	DeleteUser(c context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error)
}

type userGrpcHandler struct {
	service service.UserService
	cfg     *config.Config
	userpb.UnimplementedUserServiceServer
}

func NewUserGrpcHandler(grpc *grpc.Server, service service.UserService, cfg *config.Config) {
	handler := &userGrpcHandler{
		service: service,
		cfg:     cfg,
	}

	userpb.RegisterUserServiceServer(grpc, handler)
}

func (h *userGrpcHandler) GetAllUsers(
	c context.Context,
	req *userpb.GetAllUsersRequest,
) (*userpb.GetAllUsersResponse, error) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return nil, status.Errorf(err.Code, "%s", err.Error())
	}

	var protoUsers []*userpb.User
	for _, user := range users {
		protoUsers = append(protoUsers, mapper.MapUserEntityToProto(user))
	}

	res := &userpb.GetAllUsersResponse{
		Users: protoUsers,
	}
	return res, nil
}

func (h *userGrpcHandler) GetUserById(
	c context.Context,
	req *userpb.GetUserByIdRequest,
) (*userpb.GetUserByIdResponse, error) {
	user, err := h.service.GetUserById(uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Errorf(err.Code, "%s", err.Error())
	}

	res := &userpb.GetUserByIdResponse{
		User: mapper.MapUserEntityToProto(user),
	}
	return res, nil
}

func (h *userGrpcHandler) GetUserByEmail(
	c context.Context,
	req *userpb.GetUserByEmailRequest,
) (*userpb.GetUserByEmailResponse, error) {
	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		return nil, status.Errorf(err.Code, "%s", err.Error())
	}

	res := &userpb.GetUserByEmailResponse{
		User: mapper.MapUserEntityToProto(user),
	}
	return res, nil
}

func (h *userGrpcHandler) CreateUser(
	c context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.CreateUserResponse, error) {
	user, err := h.service.CreateUser(&domain.User{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: req.HashedPassword,
	})
	if err != nil {
		return nil, status.Errorf(err.Code, "%s", err.Error())
	}

	res := &userpb.CreateUserResponse{
		User: mapper.MapUserEntityToProto(user),
	}
	return res, nil
}

func (h *userGrpcHandler) UpdateUser(
	c context.Context,
	req *userpb.UpdateUserRequest,
) (*userpb.UpdateUserResponse, error) {
	newUserData := &domain.User{
		FullName:           req.User.FullName,
		Email:              req.User.Email,
		HashedRefreshToken: req.User.HashedRefreshToken,
	}
	userIdUUid := uuid.MustParse(req.User.Id)

	user, err := h.service.UpdateUser(userIdUUid, newUserData)
	if err != nil {
		return nil, status.Errorf(err.Code, "%s", err.Error())
	}

	res := &userpb.UpdateUserResponse{
		User: mapper.MapUserEntityToProto(user),
	}
	return res, nil
}

func (h *userGrpcHandler) DeleteUser(
	c context.Context,
	req *userpb.DeleteUserRequest,
) (*userpb.DeleteUserResponse, error) {
	user, err := h.service.DeleteUser(uuid.MustParse(req.UserId))
	if err != nil {
		return nil, status.Errorf(err.Code, "%s", err.Error())
	}

	res := &userpb.DeleteUserResponse{
		User: mapper.MapUserEntityToProto(user),
	}
	return res, nil
}
