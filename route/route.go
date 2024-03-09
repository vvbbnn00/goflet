package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"goflet/config"
	"goflet/middleware"
	"goflet/route/api"
	"goflet/route/file"
	"io"
	"time"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	// Disable the default logger
	gin.DefaultWriter = io.Discard

	router := gin.Default()
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
