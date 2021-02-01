package oauth

import (
	"net/http"

	"github.com/PS-07/go-microservices/oauth-api/src/api/domain/oauth"
	"github.com/PS-07/go-microservices/oauth-api/src/api/services"
	"github.com/PS-07/go-microservices/src/api/utils/errors"
	"github.com/gin-gonic/gin"
)

// CreateAccessToken func
func CreateAccessToken(c *gin.Context) {
	var request oauth.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusFn(), apiErr)
		return
	}
	token, err := services.OauthService.CreateAccessToken(request)
	if err != nil {
		c.JSON(err.StatusFn(), err)
		return
	}
	c.JSON(http.StatusCreated, token)
}

// GetAccessToken func
func GetAccessToken(c *gin.Context) {
	token, err := services.OauthService.GetAccessToken(c.Param("token_id"))
	if err != nil {
		c.JSON(err.StatusFn(), err)
		return
	}
	c.JSON(http.StatusOK, token)
}
