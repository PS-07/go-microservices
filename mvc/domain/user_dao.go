package domain

import (
	"fmt"
	"net/http"

	"github.com/PS-07/go-microservices/mvc/utils"
)

type userDaoInterface interface {
	GetUser(int64) (*User, *utils.ApplicationError)
}

type userDao struct{}

var (
	users = map[int64]*User{
		1: {ID: 1, FirstName: "Magnus", LastName: "Carlsen", Email: "mc@chess.com"},
		2: {ID: 2, FirstName: "Fabiano", LastName: "Caruana", Email: "fc@chess.com"},
		3: {ID: 3, FirstName: "Ding", LastName: "Liren", Email: "dl@chess.com"},
		4: {ID: 4, FirstName: "Ian", LastName: "Nepo", Email: "in@chess.com"},
		5: {ID: 5, FirstName: "Levon", LastName: "Aronian", Email: "la@chess.com"},
	}
	// UserDao var
	UserDao userDaoInterface
)

func init() {
	UserDao = &userDao{}
}

// GetUser func
func (u *userDao) GetUser(userID int64) (*User, *utils.ApplicationError) {
	if user := users[userID]; user != nil {
		return user, nil
	}
	return nil, &utils.ApplicationError{
		Message:    fmt.Sprintf("user %v does not exists", userID),
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
	}
}
