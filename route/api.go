package route

import (
	"github.com/gin-gonic/gin"
	"goflet/route/file"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	// Register the routes
	file.RegisterRoutes(router)

	return router
}
