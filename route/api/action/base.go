package action

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/vvbbnn00/goflet/storage"
	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
)

// OnConflictAction is the action to take when the file already exists
type OnConflictAction string

const (
	// OnConflictActionOverwrite is the action to overwrite the file
	OnConflictActionOverwrite OnConflictAction = "overwrite"
	// OnConflictActionAbort is the action to abort the operation
	OnConflictActionAbort OnConflictAction = "abort"
)

// CopyMoveFileRequest is the request body for the copy and move file actions
type CopyMoveFileRequest struct {
	// SourcePath is the path of the file to copy
	SourcePath string `json:"sourcePath" binding:"required"`
	// TargetPath is the path where the file will be copied
	TargetPath string `json:"targetPath" binding:"required"`
	// OnConflict is the action to take when the file already exists
	OnConflict OnConflictAction `json:"onConflict" binding:"required"`
}

// checkPath checks if the path is not empty
func checkPath(path string, c *gin.Context) (*util.Path, error) {
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return nil, errors.New("Path is required")
	}

	pathData, err := util.ParsePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	return pathData, nil
}

// preCheckForCopyMoveRoute checks the request body and the source and target paths
func preCheckForCopyMoveRoute(c *gin.Context) (*util.Path, *util.Path, bool) {
	// Get the request body
	var req CopyMoveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debugf("Error binding request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return nil, nil, false
	}

	// Check if the source and target paths are valid
	sourcePath, err := checkPath(req.SourcePath, c)
	if err != nil {
		return nil, nil, false
	}

	// Check if the source and target paths are valid
	targetPath, err := checkPath(req.TargetPath, c)
	if err != nil {
		return nil, nil, false
	}

	// Check if the source and target paths are the same
	if sourcePath.FsPath == targetPath.FsPath {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Source and target paths are the same"})
		return nil, nil, false
	}

	// Check if the source file exists
	if !storage.FileExists(sourcePath.FsPath) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source file not found"})
		return nil, nil, false
	}

	// Check if the target file exists
	if storage.FileExists(targetPath.FsPath) {
		switch req.OnConflict {
		default:
			fallthrough
		case OnConflictActionAbort:
			c.JSON(http.StatusConflict, gin.H{"error": "File already exists"})
			return nil, nil, false
		case OnConflictActionOverwrite:
			// Delete the target file
			err := storage.DeleteFile(targetPath.FsPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting target file"})
				return nil, nil, false
			}
		}
	}

	return sourcePath, targetPath, true
}
