package services

import (
	"time"

	"github.com/PS-07/go-microservices/oauth-api/src/api/domain/oauth"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
)

type oauthService struct{}

type oauthServiceInterface interface {
	CreateAccessToken(oauth.AccessTokenRequest) (*oauth.AccessToken, errors.APIError)
	GetAccessToken(accessToken string) (*oauth.AccessToken, errors.APIError)
}

// OauthService var
var OauthService oauthServiceInterface

func init() {
	OauthService = &oauthService{}
}

// CreateAccessToken func
func (s *oauthService) CreateAccessToken(request oauth.AccessTokenRequest) (*oauth.AccessToken, errors.APIError) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	user, err := oauth.GetUserByUsernameAndPassword(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	token := oauth.AccessToken{
		UserID:  user.ID,
		Expires: time.Now().UTC().Add(24 * time.Hour).Unix(),
	}
	if err := token.Save(); err != nil {
		return nil, err
	}
	return &token, nil
}

// GetAccessToken func
func (s *oauthService) GetAccessToken(accessToken string) (*oauth.AccessToken, errors.APIError) {
	token, err := oauth.GetAccessTokenByToken(accessToken)
	if err != nil {
		return nil, err
	}
	if token.IsExpired() {
		return nil, errors.NewNotFoundError("no access token found with given credentials")
	}
	return token, err
}
