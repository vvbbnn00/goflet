// Package file provides the routes for the file package
package file

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/middleware"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/storage/upload"
	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.Engine) {
	file := router.Group("/file",
		middleware.AuthChecker(),
		middleware.FilePathChecker())
	{
		// Register the routes
		file.HEAD("/*rpath", routeGetFile)
		file.GET("/*rpath", routeGetFile)
		file.PUT("/*rpath", routePutFile)
		file.POST("/*rpath", routePostFile)
		file.DELETE("/*rpath", routeDeleteFile)
	}
}

// routePutFile handler for PUT /file/*path
// @Summary      Partial File Upload
// @Description  Create an upload session with a partial file upload, supports range requests, {path} should be the relative path of the file, starting from the root directory, e.g. /file/path/to/file.txt
// @Tags         File
// @Accept       */*
// @Param        path path string true "File path"
// @Success      202  {object} string	"Accepted"
// @Failure      400  {object} string	"Bad request"
// @Failure      403  {object} string	"Directory creation not allowed"
// @Failure      500  {object} string	"Internal server error"
// @Router       /file/{path} [put]
// @Security	 Authorization
func routePutFile(c *gin.Context) {
	relativePath := c.GetString("relativePath")

	// Parse the range
	byteStart, byteEnd, _, err := util.HeaderParseRangeUpload(c.GetHeader("Content-Range"), c.GetHeader("Content-Length"))
	if err != nil {
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": err.Error()})
		return
	}

	// Get the write stream
	writeStream, err := upload.GetTempFileWriteStream(relativePath)
	if err != nil {
		errStr := err.Error()
		if errStr == "directory_creation" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Directory creation not allowed"})
			return
		}
		log.Warnf("Error getting write stream: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}
	defer func() {
		if closeErr := writeStream.Close(); closeErr != nil {
			log.Warnf("Error closing write stream: %s", closeErr.Error())
		}
	}()

	// Write the file
	body := c.Request.Body
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)

	// Seek to the start of the range
	_, err = writeStream.Seek(byteStart, io.SeekStart)
	if err != nil {
		log.Warnf("Error seeking write stream: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// Write the range to the file
	written, err := io.CopyN(writeStream, body, byteEnd-byteStart+1)
	if err != nil {
		log.Warnf("Error writing to file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}
	if written != byteEnd-byteStart+1 {
		log.Warnf("Incomplete write: expected %d bytes, wrote %d bytes", byteEnd-byteStart+1, written)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Incomplete write"})
		return
	}

	log.Debugf("Successfully written %d bytes to %s", written, relativePath)
	c.Status(http.StatusAccepted)
}

// routePostFile handler for POST /file/*path
// @Summary      Complete File Upload
// @Description  Complete an upload session with a partial file upload. You should first upload the file with a PUT request, then complete the upload with a POST request, {path} should be the relative path of the file, starting from the root directory, e.g. /file/path/to/file.txt
// @Tags         File
// @Param        path path string true "File path"
// @Success      201  {object} string	"Created"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found or upload not started"
// @Failure      409  {object} string	"File completion in progress"
// @Failure      500  {object} string	"Internal server error"
// @Router       /file/{path} [post]
// @Security	 Authorization
func routePostFile(c *gin.Context) {
	relativePath := c.GetString("relativePath")

	err := upload.CompleteFileUpload(relativePath)
	if err != nil {
		errStr := err.Error()
		if errStr == "file_uploading" {
			c.JSON(http.StatusConflict, gin.H{"error": "The file completion is in progress"})
			return
		}
		if errStr == "file_not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		log.Warnf("Error completing file upload: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error completing file upload"})
		return
	}

	c.Status(http.StatusCreated)
}

// routeDeleteFile handler for DELETE /file/*path
// @Summary      Delete File
// @Description  Delete a file by path, {path} should be the relative path of the file, starting from the root directory, e.g. /file/path/to/file.txt
// @Tags         File
// @Param        path path string true "File path"
// @Success      204  {object} string	"Deleted"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found or upload not started"
// @Failure      500  {object} string	"Internal server error"
// @Router       /file/{path} [delete]
// @Security	 Authorization
func routeDeleteFile(c *gin.Context) {
	fsPath := c.GetString("fsPath")

	err := storage.DeleteFile(fsPath)
	if err != nil {
		errStr := err.Error()
		if errStr == "file_not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		log.Warnf("Error deleting file: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting file"})
		return
	}

	c.Status(http.StatusNoContent)
}
