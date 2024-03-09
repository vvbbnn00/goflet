// Package route provides the routes for the application
package route

import (
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/middleware"
	"github.com/vvbbnn00/goflet/route/api"
	"github.com/vvbbnn00/goflet/route/file"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	if config.GofletCfg.Debug {
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = os.Stdout
	} else {
		// Disable the default logger
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	// Log the requests
	router.Use(middleware.SafeLogger())

	// Enable CORS
	corsConfig := config.GofletCfg.HTTPConfig.Cors
	if corsConfig.Enabled {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     corsConfig.Origins,
			AllowMethods:     corsConfig.Methods,
			AllowHeaders:     corsConfig.Headers,
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}))
	}

	// Register the routes
	file.RegisterRoutes(router)
	api.RegisterRoutes(router)

	return router
}
