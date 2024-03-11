// Package meta provides the routes for the meta API
package meta

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/middleware"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/util/log"
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
// @Summary      Get File Meta
// @Description  Get the file meta data, {path} should be the relative path of the file, starting from the root directory, e.g. /meta/path/to/file.txt
// @Tags         File
// @Produce      json
// @Param        path path string true "File path"
// @Success      200  {object} model.FileInfo	"OK"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found"
// @Failure      500  {object} string	"Internal server error"
// @Router       /api/meta/{path} [get]
// @Security	 Authorization
func routeGetFileMeta(c *gin.Context) {
	fsPath := c.GetString("fsPath")
	relativePath := c.GetString("relativePath")

	// Get the file info
	fileInfo, err := storage.GetFileInfo(fsPath)
	if err != nil {
		log.Debugf("Error getting file info: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Convert absolute path to relative path
	fileInfo.FilePath = relativePath

	// If windows, replace \ with /
	if filepath.Separator == '\\' {
		fileInfo.FilePath = strings.ReplaceAll(fileInfo.FilePath, "\\", "/")
	}

	c.JSON(http.StatusOK, fileInfo)
}
