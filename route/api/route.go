package api

import (
	"github.com/gin-gonic/gin"
	"goflet/middleware"
	"goflet/route/api/image"
	"goflet/route/api/meta"
	"goflet/route/api/onlyoffice"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api", middleware.AuthChecker())
	{
		onlyoffice.RegisterRoutes(api)
		meta.RegisterRoutes(api)
		image.RegisterRoutes(api)
	}
}
