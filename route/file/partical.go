package file

import (
	"github.com/vvbbnn00/goflet/config"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/vvbbnn00/goflet/storage/upload"
	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
)

// routePutUpload handler for PUT /file/*path
// @Summary      Partial File Upload
// @Description  Create an upload session with a partial file upload, supports range requests, {path} should be the relative path of the file, starting from the root directory, e.g. /file/path/to/file.txt
// @Tags         Upload
// @Accept       */*
// @Param        path path string true "File path"
// @Success      202  {object} string	"Accepted"
// @Failure      400  {object} string	"Bad request"
// @Failure      403  {object} string	"Directory creation not allowed"
// @Failure		 413  {object} string   "File too large"
// @Failure      500  {object} string	"Internal server error"
// @Router       /upload/{path} [put]
// @Security	 Authorization
func routePutUpload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.GofletCfg.FileConfig.MaxPostSize)

	relativePath := c.GetString("relativePath")

	// Parse the range
	byteStart, byteEnd, _, err := util.HeaderParseRangeUpload(c.GetHeader("Content-Range"), c.GetHeader("Content-Length"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": err.Error()})
		return
	}

	// Get the write stream
	writeStream, err := upload.GetTempFileWriteStream(relativePath)
	if err != nil {
		errStr := err.Error()
		if errStr == "directory_creation" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Directory creation not allowed"})
			return
		}
		log.Warnf("Error getting write stream: %s", errStr)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}
	defer func() {
		if closeErr := writeStream.Close(); closeErr != nil {
			log.Debugf("Error closing write stream: %s", closeErr.Error())
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// Write the range to the file
	written, err := io.CopyN(writeStream, body, byteEnd-byteStart+1)
	if err != nil {
		// Body too large
		if err.Error() == "http: request body too large" {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File too large, please use Content-Range header to upload large files"})
			return
		}
		// Other errors
		log.Warnf("Error writing to file: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}
	if written != byteEnd-byteStart+1 {
		log.Warnf("Incomplete write: expected %d bytes, wrote %d bytes", byteEnd-byteStart+1, written)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Incomplete write"})
		return
	}

	log.Debugf("Successfully written %d bytes to %s", written, relativePath)
	c.Status(http.StatusAccepted)
}

// routePostUpload handler for POST /upload/*path
// @Summary      Complete Partial File Upload
// @Description  Complete an upload session with a partial file upload. You should first upload the file with a PUT request, then complete the upload with a POST request, {path} should be the relative path of the file, starting from the root directory, e.g. /file/path/to/file.txt
// @Tags         Upload
// @Param        path path string true "File path"
// @Success      201  {object} string	"Created"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found or upload not started"
// @Failure      409  {object} string	"File completion in progress"
// @Failure      500  {object} string	"Internal server error"
// @Router       /upload/{path} [post]
// @Security	 Authorization
func routePostUpload(c *gin.Context) {
	// Complete the file upload
	relativePath := c.GetString("relativePath")
	handleCompleteFileUpload(relativePath, c)
}

// routeDeleteUpload handler for DELETE /upload/*path
// @Summary      Cancel Upload
// @Description  Cancel an upload session, {path} should be the relative path of the file, starting from the root directory, e.g. /upload/path/to/file.txt
// @Tags         Upload
// @Param        path path string true "File path"
// @Success      204  {object} string	"Deleted"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"Upload session not found"
// @Failure      500  {object} string	"Internal server error"
// @Router       /upload/{path} [delete]
// @Security	 Authorization
func routeDeleteUpload(c *gin.Context) {
	relativePath := c.GetString("relativePath")
	err := upload.RemoveTempFile(relativePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Upload session not found"})
			return
		}
		errStr := err.Error()
		log.Debugf("Error deleting file: %s", errStr)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error deleting file"})
		return
	}

	c.Status(http.StatusNoContent)
}
