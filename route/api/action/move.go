package action

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/cache"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/util/log"
)

// routeMoveFile handler for POST /action/copy
// @Summary      Move File
// @Description  Move a file from one location to another, the performance of moving is better than copying.
// @Tags         Action
// @Accept       json
// @Produce      json
// @Param        body body CopyMoveFileRequest true "Request body"
// @Success      200  {object} string	"OK"
// @Failure      400  {object} string	"Bad request"
// @Failure      404  {object} string	"File not found"
// @Failure      409  {object} string	"File exists"
// @Failure      500  {object} string	"Internal server error"
// @Router       /api/action/move [post]
// @Security	 Authorization
func routeMoveFile(c *gin.Context) {
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

	// Move the folder of the source to the target and update the metadata
	err := storage.MoveFile(sourcePath, targetPath)
	if err != nil {
		log.Debugf("Error moving file: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error moving file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File moved"})
}
