// Package middleware provides the middleware for the application
package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
)

const (
	// AuthHeader The header that contains the JWT token
	AuthHeader = "Authorization"
	// Bearer The prefix of the JWT token in the header
	Bearer = "Bearer "
	// AuthQuery The query parameter that contains the JWT token
	AuthQuery = "token"
)

// AuthChecker ensures the request is authenticated and authorized
func AuthChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			unauthorized(c, "Missing token")
			return
		}

		claims, err := parseToken(token)
		if err != nil {
			log.Debugf("Error parsing token: %s", err.Error())
			unauthorized(c, "Invalid token")
			return
		}

		if !isAuthorized(c, claims.Permissions) {
			unauthorized(c, "Unauthorized access")
			return
		}

		c.Next()
	}
}

// extractToken Extract the JWT token from the request
func extractToken(c *gin.Context) string {
	token := c.Query(AuthQuery) // Check the query parameter
	if token != "" {
		return token
	}

	token = c.GetHeader(AuthHeader) // Check the header
	if strings.HasPrefix(token, Bearer) {
		return strings.TrimPrefix(token, Bearer)
	}

	return ""
}

// parseToken Parse the JWT token
func parseToken(tokenString string) (*util.JwtClaims, error) {
	claims, err := util.ParseJwtToken(tokenString)
	if err == nil {
		return claims, nil
	}
	return nil, err
}

// queryMatch Check if the query matches the permission query
func queryMatch(query url.Values, permQuery map[string]string) bool {
	for k, v := range permQuery {
		if query.Get(k) != v {
			return false
		}
	}
	return true
}

// isAuthorized Check if the token is authorized to access the path
func isAuthorized(c *gin.Context, permissions []util.Permission) bool {
	currentPath := c.Request.URL.Path
	method := c.Request.Method

	for _, perm := range permissions {
		// Check if the path matches
		match := util.Match(currentPath, perm.Path)
		// Check if the method matches
		if match || perm.Path == "*" {
			return util.MatchMethod(method, perm.Methods) && queryMatch(c.Request.URL.Query(), perm.Query)
		}
	}

	return false
}

// unauthorized Return an unauthorized response
func unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
	c.Abort()
}
