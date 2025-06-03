package service

import (
	"context"
	"fmt"

	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gOauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
)

type OAuthGoogleService interface {
	GetGoogleLoginUrl(state string) string
	HandleGoogleCallback(code string) (*domain.User, *apperror.Error)
}

type oauthGoogleService struct {
	googleOAuthConfig *oauth2.Config
	cfg               *config.Config
}

func NewOAuthGoogleService(cfg *config.Config) OAuthGoogleService {
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
		googleOAuthConfig: googleOAuthConfig,
		cfg:               cfg,
	}
}

func (s *oauthGoogleService) GetGoogleLoginUrl(state string) string {
	return s.googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *oauthGoogleService) HandleGoogleCallback(code string) (*domain.User, *apperror.Error) {
	token, err := s.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to exchange token: %v", err.Error()))
	}

	client := s.googleOAuthConfig.Client(context.Background(), token)
	svc, err := gOauth.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to create google service: %v", err.Error()))
	}

	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("failed to get user info: %v", err.Error()))
	}

	fmt.Println(userInfo)

	return &domain.User{
		FullName: userInfo.Name,
		Email:    userInfo.Email,
	}, nil
}
