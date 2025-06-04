package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	userpb "github.com/vasapolrittideah/money-tracker-api/generated/protobuf/user"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/model"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/repository"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OAuthGoogleService interface {
	HandleGoogleCallback(code string) (*model.SignInResponse, *apperror.Error)
	GetGoogleLoginUrl(state string) string
	getGoogleUserInfo(code string) (*googleOAuth.Userinfo, *apperror.Error)
}

type oauthGoogleService struct {
	userClient        userpb.UserServiceClient
	authService       AuthService
	authRepo          repository.AuthRepository
	googleOAuthConfig *oauth2.Config
	cfg               *config.Config
}

func NewOAuthGoogleService(
	userClient userpb.UserServiceClient,
	authService AuthService,
	authRepo repository.AuthRepository,
	cfg *config.Config,
) OAuthGoogleService {
	googleOAuthConfig := &oauth2.Config{
		ClientID:     cfg.OAuthGoogle.ClientId,
		ClientSecret: cfg.OAuthGoogle.ClientSecret,
		RedirectURL:  cfg.OAuthGoogle.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &oauthGoogleService{
		userClient:        userClient,
		authService:       authService,
		authRepo:          authRepo,
		googleOAuthConfig: googleOAuthConfig,
		cfg:               cfg,
	}
}

func (s *oauthGoogleService) GetGoogleLoginUrl(state string) string {
	return s.googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *oauthGoogleService) HandleGoogleCallback(code string) (*model.SignInResponse, *apperror.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userInfo, apperr := s.getGoogleUserInfo(code)
	if apperr != nil {
		return nil, apperr
	}

	// Check if external login already exists
	externalLogin, apperr := s.authRepo.GetExternalLoginByProviderId(userInfo.Id)
	if apperr == nil {
		return s.generateJwtResponse(externalLogin.UserId)
	}

	if apperr.Code != codes.NotFound {
		return nil, apperr
	}

	// Try to get user by email
	res, err := s.userClient.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{
		Email: userInfo.Email,
	})

	var userId uuid.UUID

	if err != nil {
		st := status.Convert(err)
		if st.Code() != codes.NotFound {
			return nil, apperror.New(st.Code(), st.Err())
		}

		// User not found, create user
		_, err := s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
			FullName:       userInfo.Name,
			Email:          userInfo.Email,
			HashedPassword: "",
		})
		if err != nil {
			st := status.Convert(err)
			return nil, apperror.New(st.Code(), st.Err())
		}

		return nil, apperror.New(codes.NotFound, fmt.Errorf("user not registered"))
	} else {
		userId = uuid.MustParse(res.User.Id)
	}

	// Create external login
	_, apperr = s.authRepo.CreateExternalLogin(&domain.ExternalLogin{
		Provider:   "GOOGLE",
		ProviderId: userInfo.Id,
		UserId:     userId,
	})
	if apperr != nil {
		return nil, apperr
	}

	return s.generateJwtResponse(userId)
}

func (s *oauthGoogleService) getGoogleUserInfo(code string) (*googleOAuth.Userinfo, *apperror.Error) {
	token, err := s.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to exchange token: %v", err.Error()))
	}

	client := s.googleOAuthConfig.Client(context.Background(), token)
	svc, err := googleOAuth.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to create google service: %v", err.Error()))
	}

	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to get user info: %v", err.Error()))
	}

	return userInfo, nil
}

func (s *oauthGoogleService) generateJwtResponse(userId uuid.UUID) (*model.SignInResponse, *apperror.Error) {
	jwt, apperr := s.authService.GenerateTokens(userId)
	if apperr != nil {
		return nil, apperr
	}

	return &model.SignInResponse{Jwt: jwt}, nil
}
