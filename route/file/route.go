package file

import (
	"github.com/gin-gonic/gin"
	"goflet/middleware"
	"goflet/storage"
	"goflet/storage/upload"
	"goflet/util"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// RegisterRoutes load all the enabled routes for the application
func RegisterRoutes(router *gin.Engine) {
	file := router.Group("/file", middleware.AuthChecker(), middleware.FilePathChecker())
	{
		// Register the routes
		file.GET("/*rpath", routeGetFile)
		file.PUT("/*rpath", routePutFile)
		file.POST("/*rpath", routePostFile)
		file.DELETE("/*rpath", routeDeleteFile)
	}
}

// routeGetFile handler for GET /file/*path
func routeGetFile(c *gin.Context) {
	fsPath := c.GetString("fsPath")

	// Get the file info
	fileInfo, err := storage.GetFileInfo(fsPath)
	if err != nil {
		log.Printf("Error getting file info: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if the file has been modified
	ifModifiedSince, err := util.HeaderDateToInt64(c.GetHeader("If-Modified-Since"))
	if err == nil {
		if fileInfo.LastModified <= ifModifiedSince {
			c.Status(http.StatusNotModified)
			return
		}
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

	fileType := fileInfo.FileMeta.MimeType
	if fileType == "" {
		fileType = "application/octet-stream" // Default to application/octet-stream
	}

	// Set common headers
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.FileMeta.FileName)
	c.Header("Content-Type", fileType)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Last-Modified", util.Int64ToHeaderDate(fileInfo.LastModified))
	c.Header("X-Uploaded-At", util.Int64ToHeaderDate(fileInfo.FileMeta.UploadedAt))
	c.Header("X-Hash-Sha1", fileInfo.FileMeta.Hash.HashSha1)
	c.Header("X-Hash-Sha256", fileInfo.FileMeta.Hash.HashSha256)
	c.Header("X-Hash-Md5", fileInfo.FileMeta.Hash.HashMd5)

	byteStart, byteEnd, err := util.HeaderParseRangeDownload(c.GetHeader("Range"), fileInfo.FileSize)
	if err != nil {
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": err.Error()})
		return
	}

	// Set the Content-Range and Content-Length headers for partial content
	contentLength := byteEnd - byteStart + 1
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Header("Content-Range", "bytes "+strconv.FormatInt(byteStart, 10)+"-"+strconv.FormatInt(byteEnd, 10)+"/"+strconv.FormatInt(fileInfo.FileSize, 10))

	// Seek to the start of the range
	_, err = file.Seek(byteStart, io.SeekStart)
	if err != nil {
		log.Printf("Error seeking file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	// Set the status code
	if contentLength == fileInfo.FileSize {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusPartialContent)
	}

	// Stream the file
	_, err = io.CopyN(c.Writer, file, contentLength)
	if err != nil {
		log.Printf("Error streaming file: %s", err.Error())
		c.JSON(500, gin.H{"error": "Error reading file"})
		return
	}
}

// routePutFile handler for PUT /file/*path
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
		log.Printf("Error getting write stream: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}
	defer func() {
		if closeErr := writeStream.Close(); closeErr != nil {
			log.Printf("Error closing write stream: %s", closeErr.Error())
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
		log.Printf("Error seeking write stream: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// Write the range to the file
	written, err := io.CopyN(writeStream, body, byteEnd-byteStart+1)
	if err != nil {
		log.Printf("Error writing to file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}
	if written != byteEnd-byteStart+1 {
		log.Printf("Incomplete write: expected %d bytes, wrote %d bytes", byteEnd-byteStart+1, written)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Incomplete write"})
		return
	}

	log.Printf("Successfully written %d bytes to %s", written, relativePath)
	c.Status(http.StatusAccepted)
}

// routePostFile handler for POST /file/*path
func routePostFile(c *gin.Context) {
	relativePath := c.GetString("relativePath")

	err := upload.CompleteFileUpload(relativePath)
	if err != nil {
		errStr := err.Error()
		if errStr == "file_not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		log.Printf("Error completing file upload: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error completing file upload"})
		return
	}

	c.Status(http.StatusCreated)
}

// routeDeleteFile handler for DELETE /file/*path
func routeDeleteFile(c *gin.Context) {
	fsPath := c.GetString("fsPath")

	err := storage.DeleteFile(fsPath)
	if err != nil {
		errStr := err.Error()
		if errStr == "file_not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		log.Printf("Error deleting file: %s", errStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting file"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
