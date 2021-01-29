package repositories

// CreateRepoRequest struct
type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateRepoResponse struct
type CreateRepoResponse struct {
	ID    int64  `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}
