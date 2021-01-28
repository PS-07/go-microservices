package services

import (
	"github.com/PS-07/go-microservices/mvc/domain"
	"github.com/PS-07/go-microservices/mvc/utils"
)

type usersService struct {}

// UsersService var
var UsersService usersService

// GetUser func
func (u *usersService) GetUser(userID int64) (*domain.User, *utils.ApplicationError) {
	user, err := domain.UserDao.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
