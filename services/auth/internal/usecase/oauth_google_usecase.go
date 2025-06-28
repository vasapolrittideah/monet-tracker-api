package usecase

import (
	"context"
	"errors"
	"time"

	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	auth "github.com/vasapolrittideah/money-tracker-api/services/auth/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/grpcerror"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type oauthGoogleUsecase struct {
	userClient        userpbv1.UserServiceClient
	authUsecase       auth.AuthUsecase
	externalAuthRepo  auth.ExternalAuthRepository
	sessionRepo       auth.SessionRepository
	oauthGoogleConfig *oauth2.Config
	config            *config.Config
}

func NewOAuthGoogleUsecase(
	userClient userpbv1.UserServiceClient,
	authUsecase auth.AuthUsecase,
	externalAuthRepo auth.ExternalAuthRepository,
	sessionRepo auth.SessionRepository,
	config *config.Config,
) auth.OAuthGoogleUsecase {
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
		externalAuthRepo:  externalAuthRepo,
		sessionRepo:       sessionRepo,
		oauthGoogleConfig: oauthGoogleConfig,
		config:            config,
	}
}

func (u *oauthGoogleUsecase) GetSignInWithGoogleURL(state string) string {
	return u.oauthGoogleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (u *oauthGoogleUsecase) HandleGoogleCallback(
	ctx context.Context,
	code, userAgent, ipAddress string,
) (*auth.TokenResponse, error) {
	userInfo, err := getGoogleUserInfo(code, u.oauthGoogleConfig)
	if err != nil {
		return nil, err
	}

	// Check if Google account already linked
	externalAuth, err := u.externalAuthRepo.GetExternalAuthByProvider(ctx, "GOOGLE", userInfo.Id)
	if err == nil {
		// Google account already linked, sign in
		return u.signInWithUserID(ctx, externalAuth.UserID, userAgent, ipAddress)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to get external auth")
	}

	// Google account not linked, find or create local user
	userID, err := u.findOrCreateLocalUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	// Link Google account to local user
	_, err = u.externalAuthRepo.CreateExternalAuth(ctx, &domain.ExternalAuth{
		Provider:   "GOOGLE",
		ProviderID: userInfo.Id,
		UserID:     userID,
	})
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to create external auth")
	}

	return u.signInWithUserID(ctx, userID, userAgent, ipAddress)
}

func (u *oauthGoogleUsecase) findOrCreateLocalUser(
	ctx context.Context,
	userInfo *googleOAuth.Userinfo,
) (uint64, error) {
	res, err := u.userClient.GetUserByEmail(ctx, &userpbv1.GetUserByEmailRequest{
		Email: userInfo.Email,
	})
	if err == nil {
		return res.User.Id, nil
	}

	appErr := grpcerror.ToAppError(err).(*apperror.AppError)
	if appErr.Code != apperror.CodeNotFound {
		return 0, appErr
	}

	createdUser, err := u.userClient.CreateUser(ctx, &userpbv1.CreateUserRequest{
		FullName: userInfo.Name,
		Email:    userInfo.Email,
		Password: "", // Password is blank since it's an OAuth user
	})
	if err != nil {
		return 0, grpcerror.ToAppError(err)
	}

	return createdUser.User.Id, nil
}

func (u *oauthGoogleUsecase) signInWithUserID(
	ctx context.Context,
	userID uint64,
	userAgent, ipAddress string,
) (*auth.TokenResponse, error) {
	res, err := u.userClient.GetUserByID(ctx, &userpbv1.GetUserByIDRequest{Id: userID})
	if err != nil {
		return nil, grpcerror.ToAppError(err)
	}

	user := res.User

	session := &domain.Session{
		UserID:    user.Id,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: time.Now().Add(u.config.JWT.RefreshTokenExpiresIn),
	}
	createdSession, err := u.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to create session")
	}

	token, err := generateTokens(userID, createdSession.ID, &u.config.JWT)
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to generate tokens")
	}

	return token, nil
}

func getGoogleUserInfo(code string, config *oauth2.Config) (*googleOAuth.Userinfo, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to exchange token for google oauth")
	}

	client := config.Client(context.Background(), token)
	svc, err := googleOAuth.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to create google oauth service")
	}

	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to get user info from google oauth service")
	}

	return userInfo, nil
}
