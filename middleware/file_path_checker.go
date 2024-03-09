package middleware

import (
	"github.com/gin-gonic/gin"
	"goflet/storage"
	"goflet/util"
	"log"
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

		cleanedPath, err := util.ClarifyPath(path)
		if err != nil {
			log.Printf("Invalid path: %s, error: %s", path, err.Error())
			c.JSON(400, gin.H{"error": "Invalid path"})
			c.Abort()
			return
		}

		relativePath, err := storage.PathToRelativePath(cleanedPath)
		if err != nil {
			log.Printf("Error converting to fs path: %s", err.Error())
			c.JSON(400, gin.H{"error": "Invalid path"})
			c.Abort()
			return
		}

		fsPath, err := storage.RelativeToFsPath(relativePath)
		if err != nil {
			log.Printf("Error converting to fs path: %s", err.Error())
			c.JSON(400, gin.H{"error": "Invalid path"})
			c.Abort()
			return
		}

		// Set the cleaned path in the context
		c.Set("cleanPath", cleanedPath)
		// Set the relative path in the context
		c.Set("relativePath", relativePath)
		// Set the fs path in the context
		c.Set("fsPath", fsPath)
		c.Next()
	}
}
