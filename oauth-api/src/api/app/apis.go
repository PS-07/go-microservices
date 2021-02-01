package app

import (
	"github.com/PS-07/go-microservices/oauth-api/src/api/controllers/oauth"
	"github.com/PS-07/go-microservices/src/api/controllers/health"
)

func mapUrls() {
	router.GET("/health", health.Health)
	router.POST("/oauth/access_token", oauth.CreateAccessToken)
	router.GET("/oauth/access_token/:token_id", oauth.GetAccessToken)
}
