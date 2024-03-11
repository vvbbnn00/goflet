package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/storage/model"
	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
)

// routeGetFile handles GET and HEAD requests for /file/*path
// @Summary      File Download
// @Description  Download a file by path, supports range requests, {path} should be the relative path of the file, starting from the root directory, e.g. /file/path/to/file.txt
// @Tags         File
// @Produce      application/octet-stream
// @Success      200  {object} string	"OK"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found"
// @Failure      500  {object} string	"Internal server error"
// @Param        path path string true "File path"
// @Router       /file/{path} [get]
// @Router       /file/{path} [head]
// @Header 200,206 {string} Content-Type "application/octet-stream"
// @Header 200,206 {string} Content-Disposition "attachment; filename=file.txt"
// @Header 200,206 {string} Last-Modified "Mon, 02 Jan 2006 15:04:05 GMT"
// @Header 200,206 {string} ETag "686897696a7c876b7e"
// @Header 200,206 {string} Cache-Control "max-age=3600"
// @Header 200,206 {string} Content-Length "1024"
// @Header 206 {string} Content-Range "bytes 0-1023/2048"
// @Security	 Authorization
func routeGetFile(c *gin.Context) {
	fsPath := c.GetString("fsPath")

	// Get the file info
	fileInfo, err := storage.GetFileInfo(fsPath)
	if err != nil {
		log.Debugf("Error getting file info: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set common headers
	SetCommonHeaders(c, &fileInfo)

	// Check if the request can be responded to without reading the file
	if CanMakeFastResponse(c, &fileInfo) {
		return
	}

	// For HEAD requests, return here after setting headers
	if c.Request.Method == http.MethodHead {
		return
	}

	// Get the file reader
	file, err := storage.GetFileReader(fsPath)
	if err != nil {
		log.Warnf("Error getting file reader: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Handle range requests
	handleRangeRequests(c, file, &fileInfo)
}

// CanMakeFastResponse checks if the request can be responded to without reading the file
func CanMakeFastResponse(c *gin.Context, fileInfo *model.FileInfo) bool {
	// Check ETag header
	etag := generateETag(fileInfo)
	if ifMatch := c.GetHeader("If-Match"); ifMatch != "" && ifMatch != etag {
		c.Status(http.StatusPreconditionFailed)
		return true
	}
	if ifNoneMatch := c.GetHeader("If-None-Match"); ifNoneMatch != "" && ifNoneMatch == etag {
		c.Status(http.StatusNotModified)
		return true
	}

	// Check Last-Modified header
	lastModified := fileInfo.LastModified
	if ifModifiedSince := c.GetHeader("If-Modified-Since"); ifModifiedSince != "" &&
		util.HeaderDateToInt64(ifModifiedSince) >= lastModified {
		c.Status(http.StatusNotModified)
		return true
	}
	if ifUnmodifiedSince := c.GetHeader("If-Unmodified-Since"); ifUnmodifiedSince != "" &&
		util.HeaderDateToInt64(ifUnmodifiedSince) < lastModified {
		c.Status(http.StatusPreconditionFailed)
		return true
	}

	return false
}

// SetCommonHeaders sets common headers for the response
func SetCommonHeaders(c *gin.Context, fileInfo *model.FileInfo) {
	c.Header("Content-Type", getContentType(fileInfo))
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.FileMeta.FileName)
	c.Header("Last-Modified", util.Int64ToHeaderDate(fileInfo.LastModified))
	c.Header("ETag", generateETag(fileInfo))

	// Set the cache control header
	if config.GofletCfg.HTTPConfig.ClientCache.Enabled {
		c.Header("Cache-Control", "max-age="+strconv.Itoa(config.GofletCfg.HTTPConfig.ClientCache.MaxAge)) // Set the max age
	}
}

// handleRangeRequests handles byte range requests
func handleRangeRequests(c *gin.Context, file *os.File, fileInfo *model.FileInfo) {
	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		c.Header("Content-Length", strconv.FormatInt(fileInfo.FileSize, 10))
		c.Status(http.StatusOK)
		_, err := io.Copy(c.Writer, file)
		if err != nil {
			log.Warnf("Error copying file: %s", err.Error())
		}
		return
	}

	byteStart, byteEnd, err := util.HeaderParseRangeDownload(rangeHeader, fileInfo.FileSize)
	if err != nil {
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": err.Error()})
		return
	}
	contentLength := byteEnd - byteStart + 1

	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", byteStart, byteEnd, fileInfo.FileSize))

	if _, err := file.Seek(byteStart, io.SeekStart); err != nil {
		log.Warnf("Error seeking file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	c.Status(http.StatusPartialContent)
	_, err = io.CopyN(c.Writer, file, contentLength)
	if err != nil {
		log.Warnf("Error copying file: %s", err.Error())
	}
}

// getContentType returns the content type of the file, defaulting to "application/octet-stream"
func getContentType(fileInfo *model.FileInfo) string {
	if fileType := fileInfo.FileMeta.MimeType; fileType != "" {
		return fileType
	}
	return "application/octet-stream"
}

// generateETag generates an ETag for the file based on its metadata
func generateETag(fileInfo *model.FileInfo) string {
	return fmt.Sprintf(`"%x-%x-%s"`, fileInfo.LastModified, fileInfo.FileSize, fileInfo.FileMeta.Hash.HashSha1)
}
