package usecase

import (
	"context"

	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	auth "github.com/vasapolrittideah/money-tracker-api/services/auth/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/grpcerror"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/hashutil"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/tokenutil"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type authUsecase struct {
	userClient userpbv1.UserServiceClient
	config     *config.Config
}

func NewAuthUsecase(userClient userpbv1.UserServiceClient, config *config.Config) auth.AuthUsecase {
	return &authUsecase{
		userClient: userClient,
		config:     config,
	}
}

func (u *authUsecase) SignUp(ctx context.Context, req *auth.SignUpRequest) (*domain.User, error) {
	newUser := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	}

	created, err := u.userClient.CreateUser(ctx, &userpbv1.CreateUserRequest{
		FullName: newUser.FullName,
		Email:    newUser.Email,
		Password: newUser.Password,
	})
	if err != nil {
		return nil, grpcerror.ToAppError(err)
	}

	res, err := u.userClient.UpdateUser(ctx, &userpbv1.UpdateUserRequest{
		Id:         created.User.Id,
		Registered: wrapperspb.Bool(true),
	})
	if err != nil {
		return nil, grpcerror.ToAppError(err)
	}

	return domain.NewUserFromProto(res.User), nil
}

func (u *authUsecase) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.TokenResponse, error) {
	res, err := u.userClient.GetUserByEmail(ctx, &userpbv1.GetUserByEmailRequest{
		Email: req.Email,
	})
	if err != nil {
		return nil, grpcerror.ToAppError(err)
	}

	user := domain.NewUserFromProto(res.User)

	if ok, err := hashutil.Verify(req.Password, user.Password); err != nil || !ok {
		return nil, apperror.NewError(apperror.ErrUnauthenticated, "invalid password")
	}

	token, err := generateTokens(user.ID, &u.config.JWT)
	if err != nil {
		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}

	hashedRefreshToken, err := hashutil.Hash(token.RefreshToken)
	if err != nil {
		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}

	if _, err = u.userClient.UpdateUser(ctx, &userpbv1.UpdateUserRequest{
		Id:           user.ID,
		RefreshToken: wrapperspb.String(hashedRefreshToken),
	}); err != nil {
		return nil, grpcerror.ToAppError(err)
	}

	return token, nil
}

func generateTokens(userID uint64, jwtConfig *config.JWTConfig) (*auth.TokenResponse, error) {
	accessToken, err := tokenutil.GenerateToken(
		jwtConfig.AccessTokenExpiresIn,
		jwtConfig.AccessTokenSecretKey,
		userID,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := tokenutil.GenerateToken(
		jwtConfig.RefreshTokenExpiresIn,
		jwtConfig.RefreshTokenSecretKey,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return &auth.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
