package middleware

import (
	"github.com/gin-gonic/gin"
	"goflet/util"
	"strings"
)

const (
	// AuthHeader The header that contains the JWT token
	AuthHeader = "Authorization"
	// Bearer The prefix of the JWT token in the header
	Bearer = "Bearer "
	// AuthQuery The query parameter that contains the JWT token
	AuthQuery = "token"
)

// AuthChecker Ensures the request is authenticated
func AuthChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Request.URL.Path  // The current path
		rawQuery := c.Request.URL.RawQuery // The raw query

		if rawQuery != "" {
			currentPath += "?" + rawQuery // Append the raw query to the current path
		}

		method := c.Request.Method

		// Get the token from the query first, because the query has higher priority
		token := c.Query(AuthQuery)
		if token == "" {
			// Get the token from the header
			token = c.GetHeader(AuthHeader)
			// If the token is not empty, it must start with "Bearer"
			if !strings.HasPrefix(token, Bearer) {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			// Remove the "Bearer" prefix
			token = strings.TrimPrefix(token, Bearer)
		}

		// If the token is empty
		if token == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		//println("Token: ", token)
		claims, err := util.ParseJwtToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !util.MatchMethod(method, claims.Methods) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// Check whether the path is match
		if !util.MatchPath(currentPath, claims.Paths) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
	}
}
