package service

import (
	"encoding/gob"
	"errors"
	"goflet/cache"
	"goflet/util"
	"goflet/util/base58"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const MetaAppend = ".meta"
const FileAppend = ".file"
const FileMetaCachePrefix = "file_meta_"

type FileHash struct {
	HashSha1   string `json:"sha1"`
	HashSha256 string `json:"sha256"`
	HashMd5    string `json:"md5"`
}

type FileMeta struct {
	Hash       FileHash `json:"hash"`       // The hash of the file
	MimeType   string   `json:"mimeType"`   // The mime type of the file
	UploadedAt int64    `json:"uploadedAt"` // The time the file was uploaded
}

type FileInfo struct {
	FilePath string `json:"filePath"` // Relative path to the base file storage path
	FileName string `json:"fileName"` // The name of the file
	FileSize int64  `json:"fileSize"` // The size of the file

	LastModified int64 `json:"lastModified"` // The last modified time of the file

	FileMeta FileMeta `json:"fileMeta"` // The metadata of the file
}

// IsImage checks if the file is an image
func (f *FileInfo) IsImage() bool {
	return strings.HasPrefix(f.FileMeta.MimeType, "image/")
}

// ConvertToFsPath converts the absolute path provided to the real file system path
func ConvertToFsPath(path string) (string, error) {
	basePath := util.GetBasePath()
	if !strings.HasPrefix(path, basePath) {
		return "", errors.New("invalid path")
	}

	// Remove the base path from the path to get the relative path
	path = path[len(basePath):]
	// Replace \ with / to ensure the path is consistent
	path = filepath.ToSlash(path)
	// Separate every part of the path
	parts := strings.Split(path, "/")
	// Encode every part of the path
	for i, part := range parts {
		parts[i] = base58.Encode([]byte(part))
	}

	// Join the parts to get the real file system path
	fsPath := filepath.Join(parts...)
	fsPath = filepath.Join(basePath, fsPath)

	return fsPath, nil
}

// GetFileInfo returns the file information for the file at the provided path
func GetFileInfo(path string) (FileInfo, error) {
	fsPath, err := ConvertToFsPath(path)
	if err != nil {
		return FileInfo{}, err
	}

	fsPath += FileAppend

	file, err := os.Stat(fsPath)
	if err != nil {
		return FileInfo{}, err
	}

	name, _ := base58.Decode(strings.TrimSuffix(filepath.Base(fsPath), FileAppend))

	fileInfo := FileInfo{
		FilePath:     path,
		FileName:     string(name),
		FileSize:     file.Size(),
		LastModified: file.ModTime().Unix(),
	}

	fileMeta := GetFileMeta(path)

	return FileInfo{
		FilePath:     fileInfo.FilePath,
		FileName:     fileInfo.FileName,
		FileSize:     fileInfo.FileSize,
		LastModified: fileInfo.LastModified,
		FileMeta:     fileMeta,
	}, nil
}

// GetFileReader returns a reader for the file at the provided path, need to close the file after use
func GetFileReader(path string) (*os.File, error) {
	fsPath, err := ConvertToFsPath(path)
	if err != nil {
		return nil, err
	}

	fsPath += FileAppend

	file, err := os.OpenFile(fsPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetFileMeta returns the file metadata for the file at the provided path
func GetFileMeta(path string) FileMeta {
	fsPath, err := ConvertToFsPath(path)
	if err != nil {
		return FileMeta{}
	}

	metaFilePath := fsPath + MetaAppend

	// Check if the file metadata is cached
	c := cache.GetCache()
	cacheKey := FileMetaCachePrefix + metaFilePath

	cachedMeta, err := c.GetString(cacheKey)
	if err == nil {
		fileMeta := FileMeta{}
		err = gob.NewDecoder(strings.NewReader(cachedMeta)).Decode(&fileMeta)
		if err != nil {
			log.Printf("Error decoding meta file: %s", err.Error())
		}
		return fileMeta
	}

	metaFile, err := os.OpenFile(metaFilePath, os.O_RDONLY, 0644)

	fileMeta := FileMeta{}

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
func UpdateFileMeta(path string, fileMeta FileMeta) error {
	oldFileMeta := GetFileMeta(path)

	// Merge the old and new file metadata
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
	fsPath, err := ConvertToFsPath(path)
	if err != nil {
		return err
	}

	metaFilePath := fsPath + MetaAppend
	metaFile, err := os.OpenFile(metaFilePath, os.O_CREATE|os.O_RDWR, 0644)
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
		cacheKey := FileMetaCachePrefix + metaFilePath
		metaFileString := strings.Builder{}
		_ = gob.NewEncoder(&metaFileString).Encode(fileMeta)
		_ = c.Set(cacheKey, metaFileString.String())
	}()

	return nil
}

// DeleteFile deletes the file at the provided path
func DeleteFile(path string) error {
	fsPath, err := ConvertToFsPath(path)
	if err != nil {
		return err
	}

	fsPath += FileAppend

	// Check if the file exists
	_, err = os.Stat(fsPath)
	if err != nil {
		return errors.New("file_not_found")
	}

	err = os.Remove(fsPath)
	if err != nil {
		return err
	}

	// Delete the meta file
	metaFilePath := fsPath + MetaAppend

	_ = os.Remove(metaFilePath)

	// Delete the file metadata cache
	go func() {
		c := cache.GetCache()
		cacheKey := FileMetaCachePrefix + metaFilePath
		_ = c.Del(cacheKey)
	}()

	return nil
}
