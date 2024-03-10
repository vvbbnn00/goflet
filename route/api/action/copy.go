package action

import (
	"github.com/gin-gonic/gin"
	"github.com/vvbbnn00/goflet/cache"
	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/util/log"
	"net/http"
)

// routeCopyFile handler for POST /action/copy
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File copied"})
}
