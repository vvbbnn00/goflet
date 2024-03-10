// Package api provides the routes for the API
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vvbbnn00/goflet/route/api/action"

	"github.com/vvbbnn00/goflet/middleware"
	"github.com/vvbbnn00/goflet/route/api/image"
	"github.com/vvbbnn00/goflet/route/api/meta"
	"github.com/vvbbnn00/goflet/route/api/onlyoffice"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api", middleware.AuthChecker())
	{
		onlyoffice.RegisterRoutes(api)
		meta.RegisterRoutes(api)
		image.RegisterRoutes(api)
		action.RegisterRoutes(api)
	}
}
