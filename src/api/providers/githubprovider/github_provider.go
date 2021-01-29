package githubprovider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PS-07/go-microservices/src/api/clients/rest"
	"github.com/PS-07/go-microservices/src/api/domain/github"
)

const (
	headerAuth       = "Authorization"
	headerAuthFormat = "token %s"
	urlCreateRepo    = "https://api.github.com/user/repos"
)

func getAuthHeader(accessToken string) string {
	return fmt.Sprintf(headerAuthFormat, accessToken)
}

// CreateRepo func
func CreateRepo(accessToken string, request github.CreateRepoRequest) (*github.CreateRepoResponse, *github.ErrorResponse) {
	headers := http.Header{}
	headers.Set(headerAuth, getAuthHeader(accessToken))

	response, err := rest.Post(urlCreateRepo, request, headers)
	if err != nil {
		log.Println(fmt.Sprintf("error trying to create new repo in github: %s", err.Error()))
		return nil, &github.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	jsonBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &github.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "invalid response body",
		}
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.ErrorResponse
		if err := json.Unmarshal(jsonBytes, &errResponse); err != nil {
			return nil, &github.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "invalid json response body",
			}
		}
		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	var result github.CreateRepoResponse
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		log.Println(fmt.Sprintf("error trying to unmarshal create repo sucessful response: %s", err.Error()))
		return nil, &github.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error trying to unmarshal github create repo response",
		}
	}

	return &result, nil
}
