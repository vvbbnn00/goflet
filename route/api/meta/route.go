package meta

import (
	"github.com/gin-gonic/gin"
	"goflet/middleware"
	"goflet/service"
	"goflet/util"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.RouterGroup) {
	onlyOffice := router.Group("/meta", middleware.AuthChecker(), middleware.FilePathChecker())
	{
		// Register the routes
		onlyOffice.GET("/*rpath", routeGetFileMeta)
	}
}

// routeGetFileMeta handler for GET /meta/*path
func routeGetFileMeta(c *gin.Context) {
	basePath := util.GetBasePath()
	cleanPath := c.GetString("cleanPath")

	// Get the file info
	fileInfo, err := service.GetFileInfo(cleanPath)
	if err != nil {
		log.Printf("Error getting file info: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Convert absolute path to relative path
	fileInfo.FilePath = filepath.Join("/", strings.TrimPrefix(cleanPath, basePath))

	// If windows, replace \ with /
	if filepath.Separator == '\\' {
		fileInfo.FilePath = strings.Replace(fileInfo.FilePath, "\\", "/", -1)
	}

	c.JSON(http.StatusOK, fileInfo)
}
