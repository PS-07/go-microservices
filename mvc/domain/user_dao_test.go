package domain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserNoUserFound(t *testing.T) {
	user, err := GetUser(0)

	assert.Nil(t, user, "we were not expecting a user with id 0")
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusNotFound, err.StatusCode)
	assert.EqualValues(t, "not_found", err.Code)
	assert.EqualValues(t, "user 0 does not exists", err.Message)
}

func TestGetUserNoError(t *testing.T) {
	user, err := GetUser(1)

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t, "Magnus", user.FirstName)
	assert.EqualValues(t, "Carlsen", user.LastName)
	assert.EqualValues(t, "mc@chess.com", user.Email)
}