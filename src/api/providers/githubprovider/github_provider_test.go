package githubprovider

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/PS-07/go-microservices/src/api/clients/rest"
	"github.com/PS-07/go-microservices/src/api/domain/github"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockup()
	os.Exit(m.Run())
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "Authorization", headerAuth)
	assert.EqualValues(t, "token %s", headerAuthFormat)
	assert.EqualValues(t, "https://api.github.com/user/repos", urlCreateRepo)
}

func TestGetAuthHeader(t *testing.T) {
	header := getAuthHeader("abc123")
	assert.EqualValues(t, "token abc123", header)
}

func TestCreateRepoErrorRestclient(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Err:        errors.New("invalid restclient response"),
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid restclient response", err.Message)
}

func TestCreateRepoInvalidResponseBody(t *testing.T) {
	rest.FlushMockups()
	invalidCloser, _ := os.Open("-asf3")
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       invalidCloser,
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid response body", err.Message)
}

func TestCreateRepoInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": 1}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid json response body", err.Message)
}

func TestCreateRepoUnauthorized(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body: ioutil.NopCloser(strings.NewReader(`{
				"message": "Requires authentication",
				"documentation_url": "https://docs.github.com/rest/reference/repos#create-a-repository-for-the-authenticated-user"
			}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode)
	assert.EqualValues(t, "Requires authentication", err.Message)
}

func TestCreateRepoInvalidSuccessResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": "7"}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error trying to unmarshal github create repo response", err.Message)
}

func TestCreateRepoNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockup(rest.Mock{
		URL:        "https://api.github.com/user/repos",
		HTTPMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body: ioutil.NopCloser(strings.NewReader(`{
				"id": 123,
				"name": "golang-tutorial",
				"full_name": "PS-07/golang-tutorial"
			}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, 123, response.ID)
	assert.EqualValues(t, "golang-tutorial", response.Name)
	assert.EqualValues(t, "PS-07/golang-tutorial", response.FullName)
}
