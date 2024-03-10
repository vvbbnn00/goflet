package util

import (
	"errors"
	"github.com/vvbbnn00/goflet/util/hash"
	"github.com/vvbbnn00/goflet/util/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vvbbnn00/goflet/config"
)

var (
	// BasePath The base path for the file storage
	BasePath string
)

// init initializes the storage package
func init() {
	BasePath = GetBasePath()
	// Ensure the base path exists
	err := os.MkdirAll(BasePath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating base path: %s", err.Error())
	}
}

type Path struct {
	// The absolute path of the file in the file system
	FsPath string `json:"fsPath"`
	// The relative path of the file
	RelativePath string `json:"relativePath"`
	// Cleaned path
	CleanedPath string `json:"cleanedPath"`
}

// PathToRelativePath converts the absolute path provided to the relative path
func PathToRelativePath(path string) (string, error) {
	if !strings.HasPrefix(path, BasePath) {
		return "", errors.New("invalid path")
	}

	// Remove the base path from the path to get the relative path
	path = path[len(BasePath):]
	// Replace \ with / to ensure the path is consistent
	path = filepath.ToSlash(path)

	return path, nil
}

// RelativeToFsPath converts the relative path provided to the real file system path
func RelativeToFsPath(path string) (string, error) {
	// Get the hash of the path
	pathHash := hash.StringSha3New256(path)
	// Get Double Index of the hash
	firstIndex := pathHash[:2]
	secondIndex := pathHash[2:4]

	// Join the parts to get the real file system path
	fsPath := filepath.Join(BasePath, firstIndex, secondIndex, pathHash)

	// Add filepath separator to the end of the path to ensure it is a folder
	if !strings.HasSuffix(fsPath, string(filepath.Separator)) {
		fsPath += string(filepath.Separator)
	}

	return fsPath, nil
}

// FsPathToRelativePath converts the real file system path to the relative path
func FsPathToRelativePath(fsPath string) string {
	// Remove the base path from the path to get the relative path
	path := fsPath[len(BasePath):]
	// Replace \ with / to ensure the path is consistent
	path = filepath.ToSlash(path)

	return path
}

// ParsePath Parse the path and return the absolute and relative paths
func ParsePath(path string) (*Path, error) {
	// Ensure the path is valid
	cleanedPath, err := ClarifyPath(path)
	if err != nil {
		log.Debugf("Invalid path: %s, error: %s", path, err.Error())
		return nil, err
	}

	// Convert the path to relative path
	relativePath, err := PathToRelativePath(cleanedPath)
	if err != nil {
		log.Debugf("Error converting to fs path: %s", err.Error())
		return nil, err
	}

	// Convert the relative path to fs path
	fsPath, err := RelativeToFsPath(relativePath)
	if err != nil {
		log.Debugf("Error converting to fs path: %s", err.Error())
		return nil, err
	}

	return &Path{
		FsPath:       fsPath,
		RelativePath: relativePath,
		CleanedPath:  cleanedPath,
	}, nil
}

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
	// if strings.Contains(path, ":") {
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
