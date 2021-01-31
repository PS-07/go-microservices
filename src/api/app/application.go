package app

import (
	// "github.com/PS-07/go-microservices/src/api/log/logrus"
	"github.com/PS-07/go-microservices/src/api/log/zap"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}

// StartApp func
func StartApp() {
	// logrus.Info("about to map the APIs", "step:1", "status:pending")
	zap.Info(
		"about to map the APIs",
		zap.Field("step", 1),
		zap.Field("status", "pending"),
	)
	mapUrls()
	// logrus.Info("APIs mapped successfully", "step:2", "status:success")
	zap.Info(
		"APIs mapped successfully",
		zap.Field("step", 2),
		zap.Field("status", "success"),
	)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
