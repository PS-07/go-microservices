package services

import (
	"net/http"
	"sync"

	"github.com/PS-07/go-microservices/src/api/config"
	"github.com/PS-07/go-microservices/src/api/domain/github"
	"github.com/PS-07/go-microservices/src/api/domain/repositories"
	"github.com/PS-07/go-microservices/src/api/providers/githubprovider"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
)

type repoService struct{}

type repoServiceInterface interface {
	CreateRepo(repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError)
	CreateRepos([]repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.APIError)
}

var (
	// RepositoryService var
	RepositoryService repoServiceInterface
)

func init() {
	RepositoryService = &repoService{}
}

func (s *repoService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError) {
	if err := input.Validate(); err != nil {
		return nil, err
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

func (s *repoService) CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.APIError) {
	var wg sync.WaitGroup
	input := make(chan repositories.CreateRepoResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	go s.handleRepoResults(&wg, input, output)
	for _, req := range requests {
		wg.Add(1)
		go s.createRepoConcurrent(req, input)
	}

	wg.Wait()
	close(input)
	result := <-output

	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}
	if successCreations == 0 {
		result.StatusCode = result.Results[0].Error.StatusFn()
	} else if successCreations == len(requests) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

func (s *repoService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepoResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for incomingEvent := range input {
		repoResult := repositories.CreateRepoResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}
		results.Results = append(results.Results, repoResult)
		wg.Done()
	}
	output <- results
}

func (s *repoService) createRepoConcurrent(req repositories.CreateRepoRequest, input chan repositories.CreateRepoResult) {
	if err := req.Validate(); err != nil {
		input <- repositories.CreateRepoResult{Error: err}
		return
	}
	result, err := s.CreateRepo(req)
	if err != nil {
		input <- repositories.CreateRepoResult{Error: err}
		return
	}
	input <- repositories.CreateRepoResult{Response: result}
}
