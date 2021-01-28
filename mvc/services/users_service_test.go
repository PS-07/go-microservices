package services

import (
	"net/http"
	"testing"

	"github.com/PS-07/go-microservices/mvc/domain"
	"github.com/PS-07/go-microservices/mvc/utils"
	"github.com/stretchr/testify/assert"
)

type userDaoMock struct {}

var (
	UserDaoMock userDaoMock
	getUserFunction func(userID int64) (*domain.User, *utils.ApplicationError)
)

func (m *userDaoMock) GetUser(userID int64) (*domain.User, *utils.ApplicationError) {
	return getUserFunction(userID)
}

func init() {
	domain.UserDao = &userDaoMock{}
}

func TestGetUserNotFoundInDatabase(t *testing.T) {
	getUserFunction = func(userID int64) (*domain.User, *utils.ApplicationError) {
		return nil, &utils.ApplicationError{
			Message: "user 0 does not exists",
			StatusCode: http.StatusNotFound,
			Code: "not_found",
		}
	}

	user, err := UsersService.GetUser(0)
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, "user 0 does not exists", err.Message)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode)
	assert.EqualValues(t, "not_found", err.Code)
}

func TestGetUserNoError(t *testing.T) {
	getUserFunction = func(userID int64) (*domain.User, *utils.ApplicationError) {
		return &domain.User{
			ID: 3,
			FirstName: "Ding",
			LastName: "Liren",
			Email: "dl@chess.com",
		}, nil
	}

	user, err := UsersService.GetUser(3)
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 3, user.ID)
	assert.EqualValues(t, "Ding", user.FirstName)
	assert.EqualValues(t, "Liren", user.LastName)
	assert.EqualValues(t, "dl@chess.com", user.Email)
}