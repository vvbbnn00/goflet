package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goflet/config"
	"goflet/storage"
	"goflet/storage/model"
	"goflet/util"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// routeGetFile handles GET and HEAD requests for /file/*path
func routeGetFile(c *gin.Context) {
	fsPath := c.GetString("fsPath")

	// Get the file info
	fileInfo, err := storage.GetFileInfo(fsPath)
	if err != nil {
		log.Printf("Error getting file info: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set common headers
	setCommonHeaders(c, &fileInfo)

	// Check if the request can be responded to without reading the file
	if canMakeFastResponse(c, &fileInfo) {
		return
	}

	// For HEAD requests, return here after setting headers
	if c.Request.Method == http.MethodHead {
		return
	}

	// Get the file reader
	file, err := storage.GetFileReader(fsPath)
	if err != nil {
		log.Printf("Error getting file reader: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Handle range requests
	handleRangeRequests(c, file, &fileInfo)
}

// canMakeFastResponse checks if the request can be responded to without reading the file
func canMakeFastResponse(c *gin.Context, fileInfo *model.FileInfo) bool {
	// If-Match and If-None-Match headers
	etag := generateETag(fileInfo)
	if ifMatch := c.GetHeader("If-Match"); ifMatch != "" {
		if ifMatch != etag {
			c.Status(http.StatusPreconditionFailed)
			return true
		}
	}
	if ifNoneMatch := c.GetHeader("If-None-Match"); ifNoneMatch != "" {
		if ifNoneMatch == etag {
			c.Status(http.StatusNotModified)
			return true
		}
	}

	// If-Modified-Since and If-Unmodified-Since headers
	if ifModifiedSince := c.GetHeader("If-Modified-Since"); ifModifiedSince != "" {
		ifModifiedSinceTime, err := util.HeaderDateToInt64(ifModifiedSince)
		if err == nil && fileInfo.LastModified <= ifModifiedSinceTime {
			c.Status(http.StatusNotModified)
			return true
		}
	}

	if ifUnmodifiedSince := c.GetHeader("If-Unmodified-Since"); ifUnmodifiedSince != "" {
		ifUnmodifiedSinceTime, err := util.HeaderDateToInt64(ifUnmodifiedSince)
		if err == nil && fileInfo.LastModified > ifUnmodifiedSinceTime {
			c.Status(http.StatusPreconditionFailed)
			return true
		}
	}

	return false
}

// setCommonHeaders sets common headers for the response
func setCommonHeaders(c *gin.Context, fileInfo *model.FileInfo) {
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
	if rangeHeader := c.GetHeader("Range"); rangeHeader != "" {
		byteStart, byteEnd, err := util.HeaderParseRangeDownload(rangeHeader, fileInfo.FileSize)
		if err != nil {
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": err.Error()})
			return
		}
		contentLength := byteEnd - byteStart + 1

		etag := generateETag(fileInfo)
		lastModified := util.Int64ToHeaderDate(fileInfo.LastModified)
		if ifRange := c.GetHeader("If-Range"); ifRange != "" {
			if ifRange != etag && ifRange != lastModified {
				// The resource has been modified, so send the entire file
				byteStart = 0
				byteEnd = fileInfo.FileSize - 1
				contentLength = fileInfo.FileSize
			}
		}

		c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
		c.Header("Content-Range", "bytes "+strconv.FormatInt(byteStart, 10)+"-"+strconv.FormatInt(byteEnd, 10)+"/"+strconv.FormatInt(fileInfo.FileSize, 10))

		if _, err := file.Seek(byteStart, io.SeekStart); err != nil {
			log.Printf("Error seeking file: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
			return
		}

		c.Status(http.StatusPartialContent)
		_, err = io.CopyN(c.Writer, file, contentLength)
		if err != nil {
			return
		}
	} else {
		c.Header("Content-Length", strconv.FormatInt(fileInfo.FileSize, 10))
		c.Status(http.StatusOK)
		_, err := io.Copy(c.Writer, file)
		if err != nil {
			return
		}
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
	return fmt.Sprintf(`"%x-%x"`, fileInfo.LastModified, fileInfo.FileSize)
}
