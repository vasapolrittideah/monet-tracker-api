package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	hashedPassword, err := s.hashUserPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userProto, err := s.createNewUser(ctx, req, hashedPassword)
	if err != nil {
		return nil, err
	}

	return &model.SignUpResponse{
		User: mapper.MapUserProtoToEntity(userProto),
	}, nil
}

func (s *authService) SignIn(req *model.SignInRequest) (*model.SignInResponse, *apperror.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.getUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := s.verifyUserPassword(user.HashedPassword, req.Password); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.generateJwtTokens(user.Id)
	if err != nil {
		return nil, err
	}

	if err := s.saveRefreshToken(ctx, user, refreshToken); err != nil {
		return nil, err
	}

	return &model.SignInResponse{
		Jwt: &entity.Jwt{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *authService) hashUserPassword(password string) (string, *apperror.Error) {
	hashedPassword, err := passwordutil.HashPassword(password)
	if err != nil {
		return "", apperror.New(codes.Internal, fmt.Errorf("failed to hash password: %v", err.Error()))
	}
	return hashedPassword, nil
}

func (s *authService) createNewUser(
	ctx context.Context,
	req *model.SignUpRequest,
	hashedPassword string,
) (*userpb.User, *apperror.Error) {
	res, grpcErr := s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if grpcErr != nil {
		st := status.Convert(grpcErr)
		logger.Error("AUTH", "%s", st.Err())
		return nil, apperror.New(st.Code(), st.Err())
	}
	return res.User, nil
}

func (s *authService) getUserByEmail(ctx context.Context, email string) (*userpb.User, *apperror.Error) {
	res, grpcErr := s.userClient.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{
		Email: email,
	})
	if grpcErr != nil {
		st := status.Convert(grpcErr)
		return nil, apperror.New(st.Code(), st.Err())
	}
	return res.User, nil
}

func (s *authService) verifyUserPassword(hashedPassword, inputPassword string) *apperror.Error {
	ok, err := passwordutil.VerifyPassword(hashedPassword, inputPassword)
	if err != nil || !ok {
		return apperror.New(codes.Unauthenticated, fmt.Errorf("password is incorrect"))
	}
	return nil
}

func (s *authService) generateJwtTokens(userId string) (string, string, *apperror.Error) {
	userIdUuid := uuid.MustParse(userId)

	accessToken, err := jwtutil.GenerateJwt(
		s.cfg.Jwt.AccessTokenExpiresIn,
		s.cfg.Jwt.AccessTokenSecretKey,
		userIdUuid,
	)
	if err != nil {
		return "", "", apperror.New(codes.Internal, fmt.Errorf("failed to generate access token: %v", err.Error()))
	}

	refreshToken, err := jwtutil.GenerateJwt(
		s.cfg.Jwt.RefreshTokenExpiresIn,
		s.cfg.Jwt.RefreshTokenSecretKey,
		userIdUuid,
	)
	if err != nil {
		return "", "", apperror.New(codes.Internal, fmt.Errorf("failed to generate refresh token: %v", err.Error()))
	}

	return accessToken, refreshToken, nil
}

func (s *authService) saveRefreshToken(ctx context.Context, user *userpb.User, refreshToken string) *apperror.Error {
	hashedRefreshToken, err := jwtutil.HashRefreshToken(refreshToken)
	if err != nil {
		return apperror.New(
			codes.Internal,
			fmt.Errorf("failed to hash newly generated refresh token: %v", err.Error()),
		)
	}

	user.HashedRefreshToken = hashedRefreshToken

	if _, err = s.userClient.UpdateUser(ctx, &userpb.UpdateUserRequest{
		User: user,
	}); err != nil {
		st := status.Convert(err)
		return apperror.New(st.Code(), st.Err())
	}

	return nil
}
