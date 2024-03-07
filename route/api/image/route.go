package image

import (
	"github.com/gin-gonic/gin"
	"goflet/middleware"
	"goflet/service"
	"goflet/service/image"
	"io"
	"log"
	"net/http"
)

const (
	maxImageProcessingSize = 20 * 1024 * 1024 // 20MB
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/image", middleware.FilePathChecker())
	{
		// Register the routes
		r.GET("/*rpath", routeGetImage)
	}
}

// routeGetImage handler for GET /image/*path
func routeGetImage(c *gin.Context) {
	cleanPath := c.GetString("cleanPath")

	// Get the file info
	fileInfo, err := service.GetFileInfo(cleanPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if the file is an image
	if !fileInfo.IsImage() {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	params := image.GetProcessParamsFromQuery(c.Request.URL.Query())

	// Check if the file is too large
	if fileInfo.FileSize > maxImageProcessingSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File too large"})
		return
	}

	// Check if the file is in the cache
	cachedFile, err := image.GetFileImageReader(cleanPath, params)
	if err == nil {
		defer func() {
			_ = cachedFile.Close()
		}()
		// Set the content type
		c.Header("Content-Type", "image/"+string(params.Format))
		c.Header("X-Cache", "HIT")
		// Copy the file to the response
		_, _ = io.Copy(c.Writer, cachedFile)
		return
	}

	// Get the file read stream
	file, err := service.GetFileReader(cleanPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}
	defer func() {
		_ = file.Close()
	}()

	imageProcessed, err := image.ProcessImage(file, params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing image"})
		return
	}

	// Set the content type
	c.Header("Content-Type", "image/"+string(params.Format))
	c.Header("X-Cache", "MISS")

	// Save the file to the cache
	go func() {
		err := image.SaveFileImageCache(cleanPath, params, *imageProcessed)
		if err != nil {
			log.Printf("Error saving image cache: %s", err.Error())
		}
	}()

	// Copy the file to the response
	_, _ = io.Copy(c.Writer, imageProcessed)

}