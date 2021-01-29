package services

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/PS-07/go-microservices/src/api/clients/rest"
	"github.com/PS-07/go-microservices/src/api/domain/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockup()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidInputName(t *testing.T) {
	request := repositories.CreateRepoRequest{}
	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusFn())
	assert.EqualValues(t, "invalid repository name", err.MessageFn())
}

func TestCreateRepoErrorFromGithub(t *testing.T) {
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

	request := repositories.CreateRepoRequest{
		Name: "test-repo",
	}
	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusFn())
	assert.EqualValues(t, "Requires authentication", err.MessageFn())
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
				"name": "test-repo",
				"owner": {"login": "PS-07"}
			}`)),
		},
	})

	request := repositories.CreateRepoRequest{
		Name: "test-repo",
	}
	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 123, result.ID)
	assert.EqualValues(t, "test-repo", result.Name)
	assert.EqualValues(t, "PS-07", result.Owner)
}
