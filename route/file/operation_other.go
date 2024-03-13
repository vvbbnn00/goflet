package file

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/storage/upload"
	"github.com/vvbbnn00/goflet/util/log"
)

// routePostFile handler for POST /file/*path
// @Summary      Upload Small File
// @Description  Upload a small file using a POST request, {path} should be the relative path of the file, starting from the root directory, e.g. /file/path/to/file.txt
// @Tags         File, Upload
// @Param        path path string true "File path"
// @Accept       multipart/form-data
// @Param        file formData file true "File"
// @Success      201  {object} string	"Created"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found or upload not started"
// @Failure      409  {object} string	"File completion in progress"
// @Failure      413  {object} string	"File too large, please use PUT method to upload large files"
// @Failure      500  {object} string	"Internal server error"
// @Router       /file/{path} [post]
// @Security	 Authorization
func routePostFile(c *gin.Context) {
	// Set the request body limit
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.GofletCfg.FileConfig.MaxPostSize)
	file, err := c.FormFile("file")

	// If error is not nil and the error is "http: request body too large", return a 413 status code
	if err != nil && err.Error() == "http: request body too large" {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File too large, please use PUT method to upload large files"})
		return
	}
	if err != nil {
		log.Warnf("Error getting file: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	// If the file is not nil, handle the single file upload
	err = handleSingleFileUpload(file, c)
	if err != nil {
		return // Avoid calling handleCompleteFileUpload
	}

	// Complete the file upload
	relativePath := c.GetString("relativePath")
	handleCompleteFileUpload(relativePath, c)
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
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		log.Warnf("Error deleting file: %s", errStr)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error deleting file"})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleSingleFileUpload handles the single file upload
func handleSingleFileUpload(file *multipart.FileHeader, c *gin.Context) error {
	// Get temp file write stream
	relativePath := c.GetString("relativePath")
	writeStream, err := upload.GetTempFileWriteStream(relativePath)
	if err != nil {
		errStr := err.Error()
		if errStr == "directory_creation" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Directory creation not allowed"})
			return err
		}
		log.Warnf("Error getting write stream: %s", errStr)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return err
	}

	// Open the file
	fileReader, err := file.Open()
	if err != nil {
		log.Warnf("Error opening file: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return err
	}

	// Copy the file to the write stream
	_, err = io.Copy(writeStream, fileReader)
	if err != nil {
		log.Warnf("Error copying file: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return err
	}

	// Close the file
	_ = fileReader.Close()
	// Close the write stream
	_ = writeStream.Close()

	return nil
}

// handleCompleteFileUpload handles the completion of the file upload
func handleCompleteFileUpload(relativePath string, c *gin.Context) {
	err := upload.CompleteFileUpload(relativePath)
	if err != nil {
		errStr := err.Error()
		if errStr == "file_uploading" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "The file completion is in progress"})
			return
		}
		if errStr == "file_not_found" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "File not found or upload not started"})
			return
		}
		log.Warnf("Error completing file upload: %s", errStr)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error completing file upload"})
		return
	}

	c.Status(http.StatusCreated)
}
