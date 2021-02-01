package oauth

import (
	"github.com/PS-07/go-microservices/src/api/utils/errors"
)

const (
	queryGetUserByUsernameAndPassword = "SELECT id, username FROM users WHERE username=? AND password=?;"
)

var (
	users = map[string]*User{
		"magnus": {ID: 1, Username: "magnus"},
		"anish":  {ID: 7, Username: "anish"},
	}
)

// GetUserByUsernameAndPassword func
func GetUserByUsernameAndPassword(username string, password string) (*User, errors.APIError) {
	user := users[username]
	if user == nil {
		return nil, errors.NewNotFoundError("no user found with given credentials")
	}
	return user, nil
}
