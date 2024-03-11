package action

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/cache"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/util/log"
)

// CreateFileRequest is the request body for the create file action
type CreateFileRequest struct {
	// Path is the path where the file will be created
	Path string `json:"path" binding:"required"`
}

// routeCreateFile handler for POST /action/create
// @Summary      Create File
// @Description  Create an empty file at the specified path, if the file already exists, the operation will fail.
// @Tags         Action
// @Accept       json
// @Produce      json
// @Param        body body CreateFileRequest true "Request body"
// @Success      200  {object} string	"OK"
// @Failure      400  {object} string	"Bad request"
// @Failure      409  {object} string	"File exists"
// @Failure      500  {object} string	"Internal server error"
// @Router       /api/action/create [post]
// @Security	 Authorization
func routeCreateFile(c *gin.Context) {
	// Get the request body
	var req CreateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debugf("Error binding request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the path is valid
	pathData, err := checkPath(req.Path, c)
	if err != nil {
		return
	}

	// Check if the file already exists
	if storage.FileExists(pathData.FsPath) {
		log.Debugf("File already exists: %s", pathData.FsPath)
		c.JSON(http.StatusConflict, gin.H{"error": "File already exists"})
		return
	}

	// Lock the file, in case file upload is in progress
	ca := cache.GetCache()
	_ = ca.SetEx(storage.CachePrefix+pathData.FsPath, true, 60)

	// Unlock the file after the operation
	defer func() {
		_ = ca.Del(storage.CachePrefix + pathData.FsPath)
	}()

	// Create the file and update the metadata
	err = storage.CreateFile(pathData)
	if err != nil {
		log.Debugf("Error creating file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File created"})
}
