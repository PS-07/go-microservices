package services

import (
	"strings"

	"github.com/PS-07/go-microservices/src/api/config"
	"github.com/PS-07/go-microservices/src/api/domain/github"
	"github.com/PS-07/go-microservices/src/api/domain/repositories"
	"github.com/PS-07/go-microservices/src/api/providers/githubprovider"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
)

type repoService struct{}

type repoServiceInterface interface {
	CreateRepo(repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError)
}

var (
	// RepositoryService var
	RepositoryService repoServiceInterface
)

func init() {
	RepositoryService = &repoService{}
}

func (s *repoService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, errors.NewBadRequestError("invalid repository name")
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	response, err := githubprovider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		return nil, errors.NewAPIError(err.StatusCode, err.Message)
	}

	result := repositories.CreateRepoResponse{
		ID:    response.ID,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}
	return &result, nil
}
