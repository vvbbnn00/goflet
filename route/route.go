// Package route provides the routes for the application
package route

import (
	"io"
	"os"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/vvbbnn00/goflet/base"
	"github.com/vvbbnn00/goflet/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/middleware"
	"github.com/vvbbnn00/goflet/route/api"
	"github.com/vvbbnn00/goflet/route/file"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes() *gin.Engine {
	if config.GofletCfg.Debug {
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = os.Stdout
	} else {
		// Disable the default logger
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
	}

	// Router should be created after setting the mode
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

	// Enable swagger doc if it is enabled
	if config.GofletCfg.SwaggerDocEnabled {
		docs.SwaggerInfo.Version = base.Version
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
