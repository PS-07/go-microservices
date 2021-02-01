package oauth

import (
	"fmt"

	"github.com/PS-07/go-microservices/src/api/utils/errors"
)

var (
	tokens = make(map[string]*AccessToken, 0)
)

// Save func
func (at *AccessToken) Save() errors.APIError {
	at.AccessToken = fmt.Sprintf("USR_%d", at.UserID)
	tokens[at.AccessToken] = at
	return nil
}

// GetAccessTokenByToken func
func GetAccessTokenByToken(accessToken string) (*AccessToken, errors.APIError) {
	token := tokens[accessToken]
	if token == nil {
		return nil, errors.NewNotFoundError("no access token found with given credentials")
	}
	return token, nil
}
