// Package storage provides functions to interact with the file storage
package storage

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vvbbnn00/goflet/cache"
	"github.com/vvbbnn00/goflet/storage/model"
	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
)

const (
	maxRetryCount = 100          // The maximum number of times to retry the file uploadconst (
	CachePrefix   = "uploading:" // CachePrefix is the cache prefix for the file upload
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

// FileExists returns true if the file at the provided path exists
func FileExists(fsPath string) bool {
	_, err := os.Stat(filepath.Join(fsPath, model.FileAppend))
	return err == nil
}

// GetFileInfo returns the file information for the file at the provided path, will auto add fileAppend to the path
func GetFileInfo(fsPath string) (model.FileInfo, error) {
	filePath := filepath.Join(fsPath, model.FileAppend)
	return GetInfo(filePath)
}

// GetInfo returns the file information for the file at the provided path
func GetInfo(fsPath string) (model.FileInfo, error) {
	fi, err := os.Stat(fsPath)
	if err != nil {
		return model.FileInfo{}, err
	}

	fileInfo := model.FileInfo{
		FilePath:     fsPath,
		FileSize:     fi.Size(),
		LastModified: fi.ModTime().Unix(),
	}

	metaPath := filepath.Dir(fsPath)
	fileMeta := GetFileMeta(metaPath)

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
			log.Warnf("Error decoding meta file: %s", err.Error())
		}
		return fileMeta
	}

	log.Debugf("Cache miss: %s", metaFilePath)
	metaFile, err := os.OpenFile(metaFilePath, os.O_RDONLY, model.FilePerm)

	fileMeta := model.FileMeta{}

	if err == nil {
		gerr := gob.NewDecoder(metaFile).Decode(&fileMeta)
		if gerr != nil {
			log.Warnf("Error decoding meta file: %s", gerr.Error())
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
func UpdateFileMeta(fsPath string, fileMeta model.FileMeta) error {
	oldFileMeta := GetFileMeta(fsPath)

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
	metaFilePath := filepath.Join(fsPath, model.MetaAppend)
	tmpFilePath := filepath.Join(fsPath, "tmp-meta-"+util.RandomString(10))
	metaFile, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_RDWR, model.FilePerm)
	if err != nil {
		return err
	}
	err = gob.NewEncoder(metaFile).Encode(fileMeta)
	if err != nil {
		return err
	}
	// Close the file
	_ = metaFile.Close()

	// Move new file meta
	err = RenameFile(tmpFilePath, metaFilePath)
	if err != nil {
		_ = os.Remove(tmpFilePath)
		return err
	}

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

// RenameFile moves the file from the old path to the new path
func RenameFile(oldPath, newPath string) error {
	retryCount := 0

	// Rename the temporary file to the final file
	for {
		err := os.Rename(oldPath, newPath) // This will replace the file if it already exists
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(retryCount) * time.Millisecond) // Sleep for a while before retrying, max 5 seconds in total
		retryCount++
		if retryCount >= maxRetryCount {
			return err // Max retry count exceeded
		}
	}
}

// copyFile copies the file from the source to the destination
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Debugf("Error opening source file: %s", err.Error())
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	// Make sure the destination folder exists
	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		log.Debugf("Error creating destination folder: %s", err.Error())
		return err
	}

	// Create the destination file
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, model.FilePerm)
	if err != nil {
		log.Debugf("Error creating destination file: %s", err.Error())
		return err
	}
	defer func() {
		_ = dstFile.Close()
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		log.Debugf("Error copying file: %s", err.Error())
		return err
	}
	return nil
}

// copyFolderContents copies the contents of the folder from the source to the destination
func copyFolderContents(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		log.Debugf("Error reading directory: %s", err.Error())
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())
		if err := copyFile(srcPath, dstPath); err != nil {
			log.Debugf("Error copying file: %s", err.Error())
			return err
		}
	}
	return nil
}

// CopyFile copies the whole folder of the source to the target and update the metadata
// src and dst are absolute path of file
func CopyFile(src, dst *util.Path) error {
	// Check if the source file exists
	if !FileExists(src.FsPath) {
		return errors.New("source_file_not_found")
	}

	// Copy the folder
	err := copyFolderContents(src.FsPath, dst.FsPath)
	if err != nil {
		log.Debugf("Error copying folder contents: %s", err.Error())
		return err
	}

	log.Debugf("Successfully copied folder contents from %s to %s", src.FsPath, dst.FsPath)
	log.Debugf("Relative path: %s -> %s", src.RelativePath, dst.RelativePath)
	log.Debugf("Cleaned path: %s -> %s", src.CleanedPath, dst.CleanedPath)

	// Update the metadata
	srcMeta := GetFileMeta(src.FsPath)
	srcMeta.RelativePath = dst.RelativePath
	srcMeta.FileName = filepath.Base(dst.RelativePath)
	srcMeta.UploadedAt = time.Now().Unix()

	// Update the metadata
	return UpdateFileMeta(dst.FsPath, srcMeta)
}

// MoveFile moves the whole folder of the source to the target and update the metadata
// src and dst are absolute path of file
func MoveFile(src, dst *util.Path) error {
	// Check if the source file exists
	if !FileExists(src.FsPath) {
		return errors.New("source_file_not_found")
	}

	// Make sure the destination folder exists
	err := os.MkdirAll(filepath.Dir(dst.FsPath), os.ModePerm)
	if err != nil {
		log.Debugf("Error creating destination folder: %s", err.Error())
		return err
	}

	// Remove the destination folder
	err = os.RemoveAll(dst.FsPath)
	if err != nil {
		log.Debugf("Error removing destination folder: %s", err.Error())
		return err
	}

	// Move the folder
	err = os.Rename(src.FsPath, dst.FsPath)
	if err != nil {
		log.Debugf("Error moving folder: %s", err.Error())
		return err
	}

	// Update the metadata
	metaData := GetFileMeta(dst.FsPath)
	metaData.FileName = filepath.Base(dst.FsPath)
	metaData.RelativePath = dst.RelativePath
	err = UpdateFileMeta(dst.FsPath, metaData)

	if err != nil {
		log.Debugf("Error updating file metadata: %s", err.Error())
		return err
	}

	log.Debugf("Successfully moved folder from %s to %s", src, dst)
	return nil
}

// CreateFile creates a new file at the provided path and updates the metadata
func CreateFile(pathData *util.Path) error {
	// Make sure the folder exists
	err := os.MkdirAll(pathData.FsPath, os.ModePerm)
	if err != nil {
		log.Debugf("Error creating folder: %s", err.Error())
		return err
	}

	// Create the file
	filePath := filepath.Join(pathData.FsPath, model.FileAppend)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, model.FilePerm)
	if err != nil {
		log.Debugf("Error creating file: %s", err.Error())
		return err
	}
	// Close the file
	_ = file.Close()

	// Update the metadata
	fileMeta := model.FileMeta{
		FileName:     filepath.Base(pathData.FsPath),
		RelativePath: pathData.RelativePath,
		UploadedAt:   time.Now().Unix(),
	}

	return UpdateFileMeta(pathData.FsPath, fileMeta)
}
