package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}

// StartApp func
func StartApp() {
	mapUrls()

	fmt.Println("server listening on port 8080...")
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
