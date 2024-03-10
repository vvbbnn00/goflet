package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
)

// FilePathChecker Ensures the path is valid and does not contain any path traversal
func FilePathChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("rpath")
		// Path cannot be empty
		if path == "" {
			c.JSON(400, gin.H{"error": "Path is required"})
			c.Abort()
			return
		}

		pathData, err := util.ParsePath(path)
		if err != nil {
			log.Debugf("Invalid path: %s, error: %s", path, err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Set the cleaned path in the context
		c.Set("cleanPath", pathData.CleanedPath)
		// Set the relative path in the context
		c.Set("relativePath", pathData.RelativePath)
		// Set the fs path in the context
		c.Set("fsPath", pathData.FsPath)
		c.Next()
	}
}
