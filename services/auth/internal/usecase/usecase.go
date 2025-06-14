package usecase

import (
	"context"
	"time"

	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/hashutil"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/tokenutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type authUsercase struct {
	userClient userpbv1.UserServiceClient
	config     *config.Config
}

func NewAuthUsecase(userClient userpbv1.UserServiceClient, config *config.Config) domain.AuthUsecase {
	return &authUsercase{
		userClient: userClient,
		config:     config,
	}
}

func (u *authUsercase) SignUp(req *domain.SignUpRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newUser := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := u.userClient.CreateUser(ctx, &userpbv1.CreateUserRequest{
		FullName: newUser.FullName,
		Email:    newUser.Email,
		Password: newUser.Password,
	})
	if err != nil {
		st := status.Convert(err)
		return nil, status.Errorf(st.Code(), "failed to create user: %s", st.Message())
	}

	return domain.NewUserFromProto(res.User), nil
}

func (u *authUsercase) SignIn(req *domain.SignInRequest) (*domain.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := u.userClient.GetUserByEmail(ctx, &userpbv1.GetUserByEmailRequest{
		Email: req.Email,
	})
	if err != nil {
		st := status.Convert(err)
		return nil, status.Errorf(st.Code(), "failed to get user: %s", st.Message())
	}

	user := domain.NewUserFromProto(res.User)

	if ok, err := hashutil.Verify(req.Password, user.Password); err != nil || !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	accessToken, err := tokenutil.GenerateToken(
		u.config.JWT.AccessTokenExpiresIn,
		u.config.JWT.AccessTokenSecretKey,
		user.ID,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generate access token: %v", err)
	}

	refreshToken, err := tokenutil.GenerateToken(
		u.config.JWT.RefreshTokenExpiresIn,
		u.config.JWT.RefreshTokenSecretKey,
		user.ID,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generate refresh token: %v", err)
	}

	hashedRefreshToken, err := hashutil.Hash(refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to hash refresh token: %v", err)
	}

	if _, err = u.userClient.UpdateUser(ctx, &userpbv1.UpdateUserRequest{
		Id:           user.ID,
		RefreshToken: wrapperspb.String(hashedRefreshToken),
	}); err != nil {
		st := status.Convert(err)
		return nil, status.Errorf(st.Code(), "failed to update user: %s", st.Message())
	}

	return &domain.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
