package repositories

import (
	"strings"

	"github.com/PS-07/go-microservices/src/api/utils/errors"
)

// CreateRepoRequest struct
type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Validate func
func (req *CreateRepoRequest) Validate() errors.APIError {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return errors.NewBadRequestError("invalid repository name")
	}
	return nil
}

// CreateRepoResponse struct
type CreateRepoResponse struct {
	ID    int64  `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// CreateReposResponse struct
type CreateReposResponse struct {
	StatusCode int                `json:"status"`
	Results    []CreateRepoResult `json:"results"`
}

// CreateRepoResult struct
type CreateRepoResult struct {
	Response *CreateRepoResponse `json:"repo"`
	Error    errors.APIError     `json:"error"`
}
