// Package onlyoffice provides the routes for the OnlyOffice API
package onlyoffice

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/middleware"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/storage/upload"
	"github.com/vvbbnn00/goflet/util/log"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.RouterGroup) {
	onlyOffice := router.Group("/onlyoffice", middleware.FilePathChecker())
	{
		// Register the routes
		onlyOffice.POST("/*rpath", routeUpdateFile)
	}
}

type onlyOfficeUpdateRequest struct {
	Status int    `json:"status"` // 2 for update
	URL    string `json:"url"`    // The URL of the file
}

// routeUpdateFile handler for POST /onlyoffice/*path
// @Summary      OnlyOffice Callback
// @Description  OnlyOffice callback for file updates, when the status is 2, the file has been updated, {path} should be the relative path of the file, starting from the root directory, e.g. /onlyoffice/path/to/file.txt
// @Tags         OnlyOffice
// @Accept       json
// @Produce      json
// @Param        path path string true "File path"
// @Param        body body onlyOfficeUpdateRequest true "OnlyOffice update request"
// @Success      200  {object} string	"OK"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found"
// @Failure      500  {object} string	"Internal server error"
// @Router       /api/onlyoffice/{path} [post]
// @Security	 Authorization
func routeUpdateFile(c *gin.Context) {
	fsPath := c.GetString("fsPath")
	relativePath := c.GetString("relativePath")

	// Bind the JSON
	o := onlyOfficeUpdateRequest{}
	err := c.BindJSON(&o)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// No need to update the file
	if o.Status != 2 {
		c.JSON(http.StatusOK, gin.H{"error": 0})
		return
	}

	// Get the file info
	_, err = storage.GetFileInfo(fsPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Get the file write stream
	file, err := upload.GetTempFileWriteStream(relativePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// Download the file from the URL provided by OnlyOffice
	resp, err := http.Get(o.URL)
	if err != nil {
		log.Warnf("Error downloading file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error downloading file"})
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Write the downloaded content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Warnf("Error writing downloaded file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// Close the file
	_ = file.Close()

	// Complete the file upload
	err = upload.CompleteFileUpload(relativePath)
	if err != nil {
		errStr := err.Error()
		log.Warnf("Error completing file upload: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error completing file upload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": 0})
}
