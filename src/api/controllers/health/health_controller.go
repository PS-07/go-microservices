package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const healthStatus = "ok"

// Health func
func Health(c *gin.Context) {
	c.String(http.StatusOK, healthStatus)
}
