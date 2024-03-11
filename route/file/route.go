// Package file provides the routes for the file package
package file

import (
	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/middleware"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.Engine) {
	f := router.Group("/file",
		middleware.AuthChecker(),
		middleware.FilePathChecker())
	{
		// Register the routes for file operations
		f.HEAD("/*rpath", routeGetFile)
		f.GET("/*rpath", routeGetFile)
		f.POST("/*rpath", routePostFile)
		f.DELETE("/*rpath", routeDeleteFile)
	}

	u := router.Group("/upload",
		middleware.AuthChecker(),
		middleware.FilePathChecker())
	{
		// Register the routes for partial file upload
		u.PUT("/*rpath", routePutUpload)
		u.POST("/*rpath", routePostUpload)
		u.DELETE("/*rpath", routeDeleteUpload)
	}
}
