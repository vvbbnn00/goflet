package meta

import (
	"github.com/gin-gonic/gin"
	"goflet/middleware"
	"goflet/storage"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.RouterGroup) {
	onlyOffice := router.Group("/meta", middleware.FilePathChecker())
	{
		// Register the routes
		onlyOffice.GET("/*rpath", routeGetFileMeta)
	}
}

// routeGetFileMeta handler for GET /meta/*path
func routeGetFileMeta(c *gin.Context) {
	fsPath := c.GetString("fsPath")
	relativePath := c.GetString("relativePath")

	// Get the file info
	fileInfo, err := storage.GetFileInfo(fsPath)
	if err != nil {
		log.Printf("Error getting file info: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Convert absolute path to relative path
	fileInfo.FilePath = relativePath

	// If windows, replace \ with /
	if filepath.Separator == '\\' {
		fileInfo.FilePath = strings.Replace(fileInfo.FilePath, "\\", "/", -1)
	}

	c.JSON(http.StatusOK, fileInfo)
}
