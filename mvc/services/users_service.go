package services

import (
	"github.com/PS-07/go-microservices/mvc/domain"
	"github.com/PS-07/go-microservices/mvc/utils"
)

// GetUser func
func GetUser(userID int64) (*domain.User, *utils.ApplicationError) {
	return domain.GetUser(userID)
}
