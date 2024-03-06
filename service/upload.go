package service

import (
	"errors"
	"goflet/config"
	"goflet/util"
	"goflet/util/base58"
	"log"
	"os"
	"path/filepath"
	"time"
)

// GetTempFileWriteStream Get a write stream for the temporary file
func GetTempFileWriteStream(path string) (*os.File, error) {
	uploadPath := util.GetUploadPath()
	canCreateFolder := config.GofletCfg.FileConfig.AllowFolderCreation

	fsPath, err := convertToFsPath(path)
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

	fsPath_, err := convertToFsPath(path)
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

	// Rename the temporary file to the final file
	err = os.Rename(tmpPath, fsPath_) // This will replace the file if it already exists
	if err != nil {
		return err
	}

	// Update the file meta
	err = UpdateFileMeta(path, FileMeta{
		UploadedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	// Update the file hash
	go func() {
		fileHash := HashFile(fsPath_)
		err := UpdateFileMeta(path, FileMeta{
			Hash: fileHash,
		})
		if err != nil {
			log.Printf("Error updating file meta: %s", err.Error())
			return
		}
	}()

	return nil
}
