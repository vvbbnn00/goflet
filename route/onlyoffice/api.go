package onlyoffice

import (
	"github.com/gin-gonic/gin"
	"goflet/middleware"
	"goflet/service"
	"io"
	"log"
	"net/http"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/onlyoffice", middleware.FilePathChecker())
	{
		// Register the routes
		v1.POST("/*rpath", routeUpdateFile)
	}
}

type onlyOfficeUpdateRequest struct {
	Status int    `json:"status"` // 2 for update
	Url    string `json:"url"`    // The URL of the file
}

// routeUpdateFile handler for POST /onlyoffice/*path
func routeUpdateFile(c *gin.Context) {
	cleanPath := c.GetString("cleanPath")

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
	_, err = service.GetFileInfo(cleanPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Get the file write stream
	file, err := service.GetTempFileWriteStream(cleanPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// Download the file from the URL provided by OnlyOffice
	resp, err := http.Get(o.Url)
	if err != nil {
		log.Printf("Error downloading file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error downloading file"})
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Write the downloaded content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Error writing downloaded file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// Close the file
	_ = file.Close()

	// Complete the file upload
	err = service.CompleteFileUpload(cleanPath)
	if err != nil {
		errStr := err.Error()
		log.Printf("Error completing file upload: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error completing file upload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": 0})
}
