package service

import (
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"goflet/config"
	"goflet/util"
	"goflet/util/base58"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetTempFileWriteStream Get a write stream for the temporary file
func GetTempFileWriteStream(path string) (*os.File, error) {
	uploadPath := util.GetUploadPath()
	canCreateFolder := config.GofletCfg.FileConfig.AllowFolderCreation

	fsPath, err := ConvertToFsPath(path)
	if err != nil {
		return nil, err
	}

	// If it has subdirectory, check whether the directory can be created
	dir := filepath.Dir(fsPath)
	if dir != "." && !canCreateFolder {
		return nil, errors.New("directory_creation")
	}

	fsPath = base58.Encode([]byte(fsPath)) // Encode the path to base58 for temporary file
	tmpPath := filepath.Join(uploadPath, fsPath)

	// Ensure the directory exists
	dir = filepath.Dir(tmpPath)
	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// CompleteFileUpload Complete the file upload by renaming the temporary file to the final file
func CompleteFileUpload(path string) error {
	uploadPath := util.GetUploadPath()
	canCreateFolder := config.GofletCfg.FileConfig.AllowFolderCreation

	fsPath_, err := ConvertToFsPath(path)
	if err != nil {
		return err
	}

	fsPath := base58.Encode([]byte(fsPath_)) // Encode the path to base58 for temporary file
	tmpPath := filepath.Join(uploadPath, fsPath)

	// Check if the temporary file exists
	_, err = os.Stat(tmpPath)
	if err != nil {
		return errors.New("file_not_found")
	}

	// Ensure the directory exists
	dir := filepath.Dir(fsPath_)
	if dir != "." && !canCreateFolder {
		err := os.Remove(tmpPath)
		if err != nil {
			log.Printf("Error removing temporary file: %s", err.Error())
		}
		return errors.New("directory_creation")
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fsPath_ += FileAppend

	// Open the temporary file to get file header info
	mimeType, err := mimetype.DetectFile(tmpPath)
	if err != nil {
		return err
	}
	mimeTypeStr := mimeType.String()
	// If the file type is like html, xml, etc, set it to text/plain
	if strings.HasPrefix(mimeTypeStr, "text/") {
		mimeTypeStr = "text/plain"
	}

	// Rename the temporary file to the final file
	err = os.Rename(tmpPath, fsPath_) // This will replace the file if it already exists
	if err != nil {
		return err
	}

	// Update the file meta
	err = UpdateFileMeta(path, FileMeta{
		UploadedAt: time.Now().Unix(),
		MimeType:   mimeTypeStr,
	})
	if err != nil {
		return err
	}

	// Update the file hash
	go func() {
		HashFileAsync(fsPath_, path)
	}()

	// Remove image cache ending with .image_*
	go func() {
		err := removeImageCache(path)
		if err != nil {
			log.Printf("Error removing image cache: %s", err.Error())
		}
	}()

	return nil
}

// removeImageCache remove the image cache
func removeImageCache(path string) error {
	fsPath, err := ConvertToFsPath(path)
	if err != nil {
		return err
	}

	// Remove the file from the cache
	pathPattern := fsPath + ".image_*"
	println(err, pathPattern)
	files, err := filepath.Glob(pathPattern)

	if err != nil {
		return err
	}

	// Remove the files
	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			continue
		}
	}

	return nil
}
