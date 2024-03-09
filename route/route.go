package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"goflet/config"
	"goflet/route/api"
	"goflet/route/file"
	"time"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes() *gin.Engine {
	router := gin.Default()

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
