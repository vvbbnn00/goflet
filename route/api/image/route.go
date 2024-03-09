package image

import (
	"github.com/gin-gonic/gin"
	"goflet/config"
	"goflet/middleware"
	"goflet/route/file"
	"goflet/storage"
	"goflet/storage/image"
	"goflet/util"
	"goflet/util/log"
	"io"
	"net/http"
	"strconv"
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
	fsPath := c.GetString("fsPath")

	// Get the file info
	fileInfo, err := storage.GetFileInfo(fsPath)
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
	if fileInfo.FileSize > config.GofletCfg.ImageConfig.MaxFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File too large"})
		return
	}

	// Check if the file is in the cache
	cachedFileInfo, _ := image.GetFileImageInfo(fsPath, params)
	cachedFile, err := image.GetFileImageReader(fsPath, params)

	if err == nil && cachedFileInfo.FileSize > 0 {
		// Check if the file has been modified
		if file.CanMakeFastResponse(c, &cachedFileInfo) {
			return
		}
		defer func() {
			_ = cachedFile.Close()
		}()
		// Set the content type
		file.SetCommonHeaders(c, &cachedFileInfo)
		c.Header("Content-Disposition", "")
		c.Header("X-Cache", "HIT")
		// Copy the file to the response
		_, _ = io.Copy(c.Writer, cachedFile)
		return
	}

	// Get the file read stream
	reader, err := storage.GetFileReader(fsPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}
	defer func() {
		_ = reader.Close()
	}()

	imageProcessed, err := image.ProcessImage(reader, params)

	if err != nil {
		if err.Error() == "image size is too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Image size is too large"})
			return
		}
		log.Warnf("Error processing image: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing image"})
		return
	}

	// Set the content type
	c.Header("Content-Type", "image/"+string(params.Format))
	c.Header("Content-Length", strconv.Itoa(imageProcessed.Len()))
	c.Header("Last-Modified", util.Int64ToHeaderDate(fileInfo.LastModified))
	c.Header("X-Cache", "MISS")

	// Save the file to the cache
	go func() {
		err := image.SaveFileImageCache(fsPath, params, *imageProcessed)
		if err != nil {
			log.Warnf("Error saving image cache: %s", err.Error())
		}
	}()

	// Copy the file to the response
	_, _ = io.Copy(c.Writer, imageProcessed)

}
