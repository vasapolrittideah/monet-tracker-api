package usecase

import (
	"context"
	"errors"
	"time"

	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type oauthGoogleUsecase struct {
	userClient        userpbv1.UserServiceClient
	authUsecase       domain.AuthUsecase
	authRepo          domain.AuthRepository
	oauthGoogleConfig *oauth2.Config
	config            *config.Config
}

func NewOAuthGoogleUsecase(
	userClient userpbv1.UserServiceClient,
	authUsecase domain.AuthUsecase,
	authRepo domain.AuthRepository,
	config *config.Config,
) domain.OAuthGoogleUsecase {
	oauthGoogleConfig := &oauth2.Config{
		ClientID:     config.OAuthGoogle.ClientID,
		ClientSecret: config.OAuthGoogle.ClientSecret,
		RedirectURL:  config.OAuthGoogle.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &oauthGoogleUsecase{
		userClient:        userClient,
		authUsecase:       authUsecase,
		authRepo:          authRepo,
		oauthGoogleConfig: oauthGoogleConfig,
		config:            config,
	}
}

func (u *oauthGoogleUsecase) GetSignInWithGoogleURL(state string) string {
	return u.oauthGoogleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (u *oauthGoogleUsecase) HandleGoogleCallback(code string) (*domain.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userInfo, err := getGoogleUserInfo(code, u.oauthGoogleConfig)
	if err != nil {
		return nil, err
	}

	externalAuth, err := u.authRepo.GetExternalAuthByProviderID(userInfo.Id)
	if err == nil {
		return generateTokens(externalAuth.UserID, &u.config.JWT)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "failed to get external auth: %v", err)
	}

	res, err := u.userClient.GetUserByEmail(ctx, &userpbv1.GetUserByEmailRequest{
		Email: userInfo.Email,
	})

	var userID uint64
	if err != nil {
		st := status.Convert(err)
		if st.Code() != codes.NotFound {
			return nil, status.Error(st.Code(), st.Err().Error())
		}

		_, err := u.userClient.CreateUser(ctx, &userpbv1.CreateUserRequest{
			FullName: userInfo.Name,
			Email:    userInfo.Email,
			Password: "",
		})
		if err != nil {
			st := status.Convert(err)
			return nil, status.Error(st.Code(), st.Err().Error())
		}

		return nil, status.Error(codes.NotFound, "user not register yet")
	} else {
		userID = res.User.Id
	}

	_, err = u.authRepo.CreateExternalAuth(&domain.ExternalAuth{
		ProviderID: userInfo.Id,
		Provider:   "GOOGLE",
		UserID:     userID,
	})
	if err != nil {
		return nil, err
	}

	return generateTokens(userID, &u.config.JWT)
}

func getGoogleUserInfo(code string, config *oauth2.Config) (*googleOAuth.Userinfo, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to exchange token: %v", err)
	}

	client := config.Client(context.Background(), token)
	svc, err := googleOAuth.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create google service: %v", err)
	}

	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user info: %v", err)
	}

	return userInfo, nil
}
