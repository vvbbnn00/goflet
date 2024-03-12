// Package image provides the routes for the image API
package image

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/middleware"
	"github.com/vvbbnn00/goflet/route/file"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/storage/image"
	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
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
// @Summary      Get Image
// @Description  Get processed image, {path} should be the relative path of the file, starting from the root directory, e.g. /image/path/to/image.jpg
// @Tags         Image
// @Produce      image/jpeg, image/png, image/gif
// @Param        path path string true "File path"
// @Param        w query int false "Width"
// @Param        h query int false "Height"
// @Param        q query int false "Quality, 0-100"
// @Param        f query string false "Format" Enums(jpg, png, gif)
// @Param        a query int false "Angle, 0-360"
// @Param        s query string false "Scale type" Enums(fit, fill, resize, fit_width, fit_height)
// @Success      200  {object} string	"OK"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found"
// @Failure      413  {object} string	"File too large"
// @Failure      500  {object} string	"Internal server error"
// @Router       /api/image/{path} [get]
// @Security	 Authorization
func routeGetImage(c *gin.Context) {
	fsPath := c.GetString("fsPath")

	// Get the file info
	fileInfo, err := storage.GetFileInfo(fsPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if the file is an image
	if !fileInfo.IsImage() {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	params := image.GetProcessParamsFromQuery(c.Request.URL.Query())

	// Check if the file is too large
	if fileInfo.FileSize > config.GofletCfg.ImageConfig.MaxFileSize {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File too large"})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}
	defer func() {
		_ = reader.Close()
	}()

	imageProcessed, err := image.ProcessImage(reader, params)

	if err != nil {
		if err.Error() == "image size is too large" {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Image size is too large"})
			return
		}
		log.Warnf("Error processing image: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error processing image"})
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
