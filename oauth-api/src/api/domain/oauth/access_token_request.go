package oauth

import (
	"strings"

	"github.com/PS-07/go-microservices/src/api/utils/errors"
)

// AccessTokenRequest struct
type AccessTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate func
func (r *AccessTokenRequest) Validate() errors.APIError {
	r.Username = strings.TrimSpace(r.Username)
	if r.Username == "" {
		return errors.NewBadRequestError("invalid username")
	}
	r.Password = strings.TrimSpace(r.Password)
	if r.Password == "" {
		return errors.NewBadRequestError("invalid username")
	}
	return nil
}
