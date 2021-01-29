package repositories

import (
	"net/http"

	"github.com/PS-07/go-microservices/src/api/domain/repositories"
	"github.com/PS-07/go-microservices/src/api/services"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
	"github.com/gin-gonic/gin"
)

// CreateRepo func
func CreateRepo(c *gin.Context) {
	var request repositories.CreateRepoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusFn(), apiErr)
		return
	}

	result, err := services.RepositoryService.CreateRepo(request)
	if err != nil {
		c.JSON(err.StatusFn(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
