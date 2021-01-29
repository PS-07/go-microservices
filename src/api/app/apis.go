package app

import (
	"github.com/PS-07/go-microservices/src/api/controllers/health"
	"github.com/PS-07/go-microservices/src/api/controllers/repositories"
)

func mapUrls() {
	router.GET("/health", health.Health)
	router.POST("/repositories", repositories.CreateRepo)
}
