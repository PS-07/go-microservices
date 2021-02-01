package repositoriesmock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	repocontroller "github.com/PS-07/go-microservices/src/api/controllers/repositories"
	"github.com/PS-07/go-microservices/src/api/domain/repositories"
	"github.com/PS-07/go-microservices/src/api/services"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
	"github.com/PS-07/go-microservices/src/api/utils/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	funcCreateRepo  func(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError)
	funcCreateRepos func(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.APIError)
)

type repoServiceMock struct{}

func (s *repoServiceMock) CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError) {
	return funcCreateRepo(request)
}

func (s *repoServiceMock) CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.APIError) {
	return funcCreateRepos(requests)
}

func TestCreateRepoNoErrorMockingService(t *testing.T) {
	services.RepositoryService = &repoServiceMock{}

	funcCreateRepo = func(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError) {
		return &repositories.CreateRepoResponse{
			ID:    123,
			Name:  "test-repo",
			Owner: "PS-07",
		}, nil
	}

	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "test-repo"}`))
	response := httptest.NewRecorder()
	c := testutils.GetMockedContext(request, response)

	repocontroller.CreateRepo(c)
	assert.EqualValues(t, http.StatusCreated, response.Code)

	var result repositories.CreateRepoResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, result.ID)
	assert.EqualValues(t, "test-repo", result.Name)
	assert.EqualValues(t, "PS-07", result.Owner)
}
