package service

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/vasapolrittideah/money-tracker-api/generated/protobuf/user"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/model"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain/entity"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
	"github.com/vasapolrittideah/money-tracker-api/shared/mapper"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/jwtutil"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/passwordutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	SignUp(req *model.SignUpRequest) (*model.SignUpResponse, *apperror.Error)
	SignIn(req *model.SignInRequest) (*model.SignInResponse, *apperror.Error)
}

type authService struct {
	userClient userpb.UserServiceClient
	cfg        *config.Config
}

func NewAuthService(userClient userpb.UserServiceClient, cfg *config.Config) AuthService {
	return &authService{
		userClient: userClient,
		cfg:        cfg,
	}
}

func (s *authService) SignUp(req *model.SignUpRequest) (*model.SignUpResponse, *apperror.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := passwordutil.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := entity.User{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	res, grpcErr := s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		FullName:       newUser.FullName,
		Email:          newUser.Email,
		HashedPassword: newUser.HashedPassword,
	})
	if grpcErr != nil {
		st := status.Convert(grpcErr)
		logger.Error("AUTH", "%s", st.Err())
		return nil, apperror.New(st.Code(), st.Err())
	}

	return &model.SignUpResponse{
		User: mapper.MapUserProtoToEntity(res.User),
	}, nil
}

func (s *authService) SignIn(req *model.SignInRequest) (*model.SignInResponse, *apperror.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, grpcErr := s.userClient.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{
		Email: req.Email,
	})
	if grpcErr != nil {
		st := status.Convert(grpcErr)
		return nil, apperror.New(st.Code(), st.Err())
	}

	user := mapper.MapUserProtoToEntity(res.User)

	if ok, err := passwordutil.VerifyPassword(user.HashedPassword, req.Password); err != nil || !ok {
		return nil, apperror.New(codes.Unauthenticated, fmt.Errorf("password is incorrect"))
	}

	accessToken, err := jwtutil.GenerateJwt(
		s.cfg.Jwt.AccessTokenExpiresIn,
		s.cfg.Jwt.AccessTokenSecretKey,
		user.Id,
	)
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to generate access token: %v", err.Error()))
	}

	refreshToken, err := jwtutil.GenerateJwt(
		s.cfg.Jwt.RefreshTokenExpiresIn,
		s.cfg.Jwt.RefreshTokenSecretKey,
		user.Id,
	)
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to generate refresh token: %v", err.Error()))
	}

	hashedRefreshToken, err := jwtutil.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, apperror.New(
			codes.Internal,
			fmt.Errorf("failed to hash newly generated refresh token: %v", err.Error()),
		)
	}

	if _, err = s.userClient.UpdateUser(ctx, &userpb.UpdateUserRequest{
		User: &userpb.User{
			HashedRefreshToken: hashedRefreshToken,
		},
	}); err != nil {
		st := status.Convert(err)
		return nil, apperror.New(st.Code(), st.Err())
	}

	jwtRes := &model.SignInResponse{
		Jwt: &entity.Jwt{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	return jwtRes, nil
}
