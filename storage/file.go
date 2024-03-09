package storage

import (
	"encoding/gob"
	"errors"
	"goflet/cache"
	"goflet/storage/model"
	"goflet/util"
	"goflet/util/hash"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	basePath string
)

// init initializes the storage package
func init() {
	basePath = util.GetBasePath()
	// Ensure the base path exists
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating base path: %s", err.Error())
	}
}

// PathToRelativePath converts the absolute path provided to the relative path
func PathToRelativePath(path string) (string, error) {
	if !strings.HasPrefix(path, basePath) {
		return "", errors.New("invalid path")
	}

	// Remove the base path from the path to get the relative path
	path = path[len(basePath):]
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
	fsPath := filepath.Join(basePath, firstIndex, secondIndex, pathHash)

	// Add filepath separator to the end of the path to ensure it is a folder
	if !strings.HasSuffix(fsPath, string(filepath.Separator)) {
		fsPath += string(filepath.Separator)
	}

	return fsPath, nil
}

// PathToFsPath converts the absolute path provided to the real file system path
func PathToFsPath(path string) (string, error) {
	relativePath, err := PathToRelativePath(path)

	if err != nil {
		return "", err
	}

	return RelativeToFsPath(relativePath)
}

// GetFileInfo returns the file information for the file at the provided path
func GetFileInfo(fsPath string) (model.FileInfo, error) {
	filePath := filepath.Join(fsPath, model.FileAppend)
	file, err := os.Stat(filePath)
	if err != nil {
		return model.FileInfo{}, err
	}

	fileInfo := model.FileInfo{
		FilePath:     fsPath,
		FileSize:     file.Size(),
		LastModified: file.ModTime().Unix(),
	}

	fileMeta := GetFileMeta(fsPath)

	return model.FileInfo{
		FilePath:     fileInfo.FilePath,
		FileSize:     fileInfo.FileSize,
		LastModified: fileInfo.LastModified,
		FileMeta:     fileMeta,
	}, nil
}

// GetFileReader returns a reader for the file at the provided path, need to close the file after use
func GetFileReader(fsPath string) (*os.File, error) {
	filePath := filepath.Join(fsPath, model.FileAppend)

	file, err := os.OpenFile(filePath, os.O_RDONLY, model.FilePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetFileMeta returns the file metadata for the file at the provided path
func GetFileMeta(fsPath string) model.FileMeta {
	metaFilePath := filepath.Join(fsPath, model.MetaAppend)

	// Check if the file metadata is cached
	c := cache.GetCache()
	cacheKey := model.FileMetaCachePrefix + metaFilePath

	cachedMeta, err := c.GetString(cacheKey)
	if err == nil {
		fileMeta := model.FileMeta{}
		err = gob.NewDecoder(strings.NewReader(cachedMeta)).Decode(&fileMeta)
		if err != nil {
			log.Printf("Error decoding meta file: %s", err.Error())
		}
		return fileMeta
	}

	metaFile, err := os.OpenFile(metaFilePath, os.O_RDONLY, model.FilePerm)

	fileMeta := model.FileMeta{}

	if err == nil {
		gerr := gob.NewDecoder(metaFile).Decode(&fileMeta)
		if gerr != nil {
			log.Printf("Error decoding meta file: %s", gerr.Error())
		}
	}

	// Close the file
	_ = metaFile.Close()

	// Cache the file metadata
	go func() {
		metaFileString := strings.Builder{}
		_ = gob.NewEncoder(&metaFileString).Encode(fileMeta)
		_ = c.Set(cacheKey, metaFileString.String())
	}()

	return fileMeta
}

// UpdateFileMeta updates the file metadata for the file at the provided path
func UpdateFileMeta(faPath string, fileMeta model.FileMeta) error {
	oldFileMeta := GetFileMeta(faPath)

	// Merge the old and new file metadata
	if fileMeta.RelativePath == "" {
		fileMeta.RelativePath = oldFileMeta.RelativePath
	}
	if fileMeta.FileName == "" {
		fileMeta.FileName = oldFileMeta.FileName
	}
	if fileMeta.UploadedAt == 0 {
		fileMeta.UploadedAt = oldFileMeta.UploadedAt
	}
	if fileMeta.Hash.HashMd5 == "" {
		fileMeta.Hash.HashMd5 = oldFileMeta.Hash.HashMd5
	}
	if fileMeta.Hash.HashSha1 == "" {
		fileMeta.Hash.HashSha1 = oldFileMeta.Hash.HashSha1
	}
	if fileMeta.Hash.HashSha256 == "" {
		fileMeta.Hash.HashSha256 = oldFileMeta.Hash.HashSha256
	}
	if fileMeta.MimeType == "" {
		fileMeta.MimeType = oldFileMeta.MimeType
	}

	// Save the new file metadata
	metaFilePath := filepath.Join(faPath, model.MetaAppend)
	metaFile, err := os.OpenFile(metaFilePath, os.O_CREATE|os.O_RDWR, model.FilePerm)
	if err != nil {
		return err
	}

	err = gob.NewEncoder(metaFile).Encode(fileMeta)
	if err != nil {
		return err
	}

	// Close the file
	_ = metaFile.Close()

	// Cache the file metadata
	go func() {
		c := cache.GetCache()
		cacheKey := model.FileMetaCachePrefix + metaFilePath
		metaFileString := strings.Builder{}
		_ = gob.NewEncoder(&metaFileString).Encode(fileMeta)
		_ = c.Set(cacheKey, metaFileString.String())
	}()

	return nil
}

// DeleteFile deletes the file at the provided path
func DeleteFile(fsPath string) error {
	// Check if the folder exists
	_, err := os.Stat(fsPath)
	if err != nil {
		return errors.New("file_not_found")
	}

	// Delete the folder and its contents
	err = os.RemoveAll(fsPath)
	if err != nil {
		return err
	}
	return nil
}