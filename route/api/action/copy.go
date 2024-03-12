package action

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/cache"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/util/log"
)

// routeCopyFile handler for POST /action/copy
// @Summary      Copy File
// @Description  Copy a file from one location to another, if you want to move a file, use the move action instead.
// @Tags         Action
// @Accept       json
// @Produce      json
// @Param        body body CopyMoveFileRequest true "Request body"
// @Success      200  {object} string	"OK"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found"
// @Failure      409  {object} string	"File exists"
// @Failure      500  {object} string	"Internal server error"
// @Router       /api/action/copy [post]
// @Security	 Authorization
func routeCopyFile(c *gin.Context) {
	sourcePath, targetPath, ok := preCheckForCopyMoveRoute(c)
	if !ok {
		return
	}

	// Lock the file, in case file upload is in progress
	ca := cache.GetCache()
	_ = ca.SetEx(storage.CachePrefix+targetPath.FsPath, true, 60)

	// Unlock the file after the operation
	defer func() {
		_ = ca.Del(storage.CachePrefix + targetPath.FsPath)
	}()

	// Copy the whole folder of the source to the target and update the metadata
	err := storage.CopyFile(sourcePath, targetPath)
	if err != nil {
		log.Debugf("Error copying file: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error copying file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File copied"})
}
