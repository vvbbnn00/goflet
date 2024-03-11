// Package action provides the routes for the action API
package action

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/action")
	{
		// Register the routes
		r.POST("/copy", routeCopyFile)
		r.POST("/move", routeMoveFile)
		r.POST("/create", routeCreateFile)
	}
}
