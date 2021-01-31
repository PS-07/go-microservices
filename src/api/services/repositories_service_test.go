package services

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/PS-07/go-microservices/src/api/clients/rest"
	"github.com/PS-07/go-microservices/src/api/domain/repositories"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
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
	request := repositories.CreateRepoRequest{Name: "test-repo"}
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
	request := repositories.CreateRepoRequest{Name: "test-repo"}
	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 123, result.ID)
	assert.EqualValues(t, "test-repo", result.Name)
	assert.EqualValues(t, "PS-07", result.Owner)
}

func TestCreateRepoConcurrentInvalidRequest(t *testing.T) {
	service := repoService{}
	request := repositories.CreateRepoRequest{}
	input := make(chan repositories.CreateRepoResult)
	
	go service.createRepoConcurrent(request, input)
	result := <- input

	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Error.StatusFn())
	assert.EqualValues(t, "invalid repository name", result.Error.MessageFn())
}

func TestCreateRepoConcurrentErrorFromGithub(t *testing.T) {
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
	service := repoService{}
	request := repositories.CreateRepoRequest{Name: "test-repo"}
	input := make(chan repositories.CreateRepoResult)
	
	go service.createRepoConcurrent(request, input)
	result := <- input

	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusUnauthorized, result.Error.StatusFn())
	assert.EqualValues(t, "Requires authentication", result.Error.MessageFn())
}

func TestCreateRepoConcurrentNoError(t *testing.T) {
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
	service := repoService{}
	request := repositories.CreateRepoRequest{Name: "test-repo"}
	input := make(chan repositories.CreateRepoResult)
	
	go service.createRepoConcurrent(request, input)
	result := <- input

	assert.NotNil(t, result)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Response)
	assert.EqualValues(t, 123, result.Response.ID)
	assert.EqualValues(t, "test-repo", result.Response.Name)
	assert.EqualValues(t, "PS-07", result.Response.Owner)
}

func TestHandleRepoResults(t *testing.T) {
	var wg sync.WaitGroup
	input := make(chan repositories.CreateRepoResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	service := repoService{}
	go service.handleRepoResults(&wg, input, output)
	
	wg.Add(1)
	go func() {
		input <- repositories.CreateRepoResult{
			Error: errors.NewBadRequestError("invalid repository name"),
		}
	} ()

	wg.Wait()
	close(input)
	result := <- output

	assert.NotNil(t, result)
	assert.EqualValues(t, 0, result.StatusCode)
	assert.EqualValues(t, 1, len(result.Results))
	assert.NotNil(t, result.Results[0].Error)	
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.StatusFn())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.MessageFn())
}

func TestCreateReposInvalidRequests(t *testing.T) {
	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "   "},
	}
	result, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusBadRequest, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	assert.Nil(t, result.Results[0].Response)
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.StatusFn())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.MessageFn())

	assert.Nil(t, result.Results[1].Response)
	assert.EqualValues(t, http.StatusBadRequest, result.Results[1].Error.StatusFn())
	assert.EqualValues(t, "invalid repository name", result.Results[1].Error.MessageFn())
}

func TestCreateReposOneSuccessOneFailure(t *testing.T) {
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
	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "test-repo"},
	}
	result, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusPartialContent, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	for _, result := range result.Results {
		if result.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest, result.Error.StatusFn())
			assert.EqualValues(t, "invalid repository name", result.Error.MessageFn())
			continue
		}
		assert.EqualValues(t, 123, result.Response.ID)
		assert.EqualValues(t, "test-repo", result.Response.Name)
		assert.EqualValues(t, "PS-07", result.Response.Owner)
	}
}

func TestCreateReposRepoAlreadyExistsFailure(t *testing.T) {
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

	requests := []repositories.CreateRepoRequest{
		{Name: "test-repo"},
		{Name: "test-repo"},
	}

	result, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusPartialContent, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	for _, result := range result.Results {
		if result.Error != nil {
			assert.EqualValues(t, http.StatusInternalServerError, result.Error.StatusFn())
			assert.EqualValues(t, "error trying to unmarshal github create repo response", result.Error.MessageFn())
			continue
		}

		assert.EqualValues(t, 123, result.Response.ID)
		assert.EqualValues(t, "test-repo", result.Response.Name)
		assert.EqualValues(t, "PS-07", result.Response.Owner)
	}
}