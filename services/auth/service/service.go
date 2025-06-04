package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	userpb "github.com/vasapolrittideah/money-tracker-api/generated/protobuf/user"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/model"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
	"github.com/vasapolrittideah/money-tracker-api/shared/mapper"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/passwordutil"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/tokenutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	SignUp(req *model.SignUpRequest) (*model.SignUpResponse, *apperror.Error)
	SignIn(req *model.SignInRequest) (*model.SignInResponse, *apperror.Error)
	GenerateTokens(userId uuid.UUID) (*domain.Jwt, *apperror.Error)
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

	hashedPassword, apperr := passwordutil.HashPassword(req.Password)
	if apperr != nil {
		return nil, apperr
	}

	newUser := domain.User{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	res, err := s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		FullName:       newUser.FullName,
		Email:          newUser.Email,
		HashedPassword: newUser.HashedPassword,
	})
	if err != nil {
		st := status.Convert(err)
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

	res, err := s.userClient.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{
		Email: req.Email,
	})
	if err != nil {
		st := status.Convert(err)
		return nil, apperror.New(st.Code(), st.Err())
	}

	user := res.User
	userIdUuid := uuid.MustParse(user.Id)

	if ok, err := passwordutil.VerifyPassword(user.HashedPassword, req.Password); err != nil || !ok {
		return nil, apperror.New(codes.Unauthenticated, fmt.Errorf("password is incorrect"))
	}

	jwt, apperr := s.GenerateTokens(userIdUuid)
	if apperr != nil {
		return nil, apperr
	}

	hashedRefreshToken, err := tokenutil.HashRefreshToken(jwt.RefreshToken)
	if err != nil {
		return nil, apperror.New(
			codes.Internal,
			fmt.Errorf("failed to hash newly generated refresh token: %v", err.Error()),
		)
	}

	user.HashedRefreshToken = hashedRefreshToken
	if _, err = s.userClient.UpdateUser(ctx, &userpb.UpdateUserRequest{
		User: user,
	}); err != nil {
		st := status.Convert(err)
		return nil, apperror.New(st.Code(), st.Err())
	}

	jwtRes := &model.SignInResponse{Jwt: jwt}

	return jwtRes, nil
}

func (s *authService) GenerateTokens(userId uuid.UUID) (*domain.Jwt, *apperror.Error) {
	accessToken, err := tokenutil.GenerateToken(
		s.cfg.Jwt.AccessTokenExpiresIn,
		s.cfg.Jwt.AccessTokenSecretKey,
		userId,
	)
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to generate access token: %v", err.Error()))
	}

	refreshToken, err := tokenutil.GenerateToken(
		s.cfg.Jwt.RefreshTokenExpiresIn,
		s.cfg.Jwt.RefreshTokenSecretKey,
		userId,
	)
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to generate refresh token: %v", err.Error()))
	}

	return &domain.Jwt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
