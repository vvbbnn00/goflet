package util

import (
	"errors"
	"goflet/config"
	"os"
	"path/filepath"
	"strings"
)

// GetPath Get the absolute path for the file storage
func GetPath(path string) string {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	// Ensure the base path ends with a separator
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}

	// Check if the base path exists
	if _, err := filepath.Abs(path); err != nil {
		// Create the base path
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return path
}

// GetBasePath Get the base path for the file storage
func GetBasePath() string {
	basePath := config.GofletCfg.FileConfig.BaseFileStoragePath
	return GetPath(basePath)
}

// GetUploadPath Get the upload path for the file storage
func GetUploadPath() string {
	uploadPath := config.GofletCfg.FileConfig.UploadPath
	return GetPath(uploadPath)
}

// ClarifyPath Ensures the path is valid and does not contain any path traversal
func ClarifyPath(path string) (string, error) {
	basePath := GetBasePath()

	// Ensure the base path is absolute
	if !filepath.IsAbs(basePath) {
		basePath, _ = filepath.Abs(basePath)
	}

	//// Check if the path contains unsupported characters
	//if strings.Contains(path, ":") {
	//	return "", errors.New("path contains unsupported characters")
	//}

	// Combine the base path and the path
	combinedPath := filepath.Join(basePath, path)
	// Clean the path
	cleanPath := filepath.Clean(combinedPath)

	// Avoid path traversal
	basePathWithSlash := basePath
	if !strings.HasSuffix(basePathWithSlash, string(filepath.Separator)) {
		basePathWithSlash += string(filepath.Separator)
	}
	if !strings.HasPrefix(cleanPath, basePathWithSlash) {
		return "", errors.New("path traversal detected")
	}

	return cleanPath, nil
}

// Match the pattern with the name
func Match(pattern, name string) bool {
	// If it has no *, just compare the strings
	if !strings.Contains(pattern, "*") {
		return pattern == name
	}

	// Split the pattern by *
	parts := strings.Split(pattern, "*")
	return strings.HasPrefix(name, parts[0]) // Only compare the first part
}

// MatchMethod Match the method with the methods
func MatchMethod(method string, methods []string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}
