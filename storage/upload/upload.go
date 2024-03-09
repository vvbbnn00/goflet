package upload

import (
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"goflet/config"
	"goflet/storage"
	"goflet/storage/hasher"
	"goflet/storage/image"
	"goflet/storage/model"
	"goflet/util"
	"goflet/util/hash"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	uploadPath      string
	canCreateFolder bool
)

// init initializes the upload package
func init() {
	uploadPath = util.GetUploadPath()
	canCreateFolder = config.GofletCfg.FileConfig.AllowFolderCreation

	// Ensure the upload path exists
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

// GetTempFileWriteStream Get a write stream for the temporary file
func GetTempFileWriteStream(relativePath string) (*os.File, error) {
	// If it has subdirectory, check whether the directory can be created
	dir := filepath.Dir(relativePath)
	if dir != "." && !canCreateFolder {
		return nil, errors.New("directory_creation")
	}

	fileName := hash.StringSha3New256(relativePath) // Get the hash of the path
	tmpPath := filepath.Join(uploadPath, fileName)

	// Ensure the directory exists
	dir = filepath.Dir(tmpPath)
	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR, model.FilePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// CompleteFileUpload Complete the file upload by renaming the temporary file to the final file
func CompleteFileUpload(relativePath string) error {
	fileName := hash.StringSha3New256(relativePath) // Get the hash of the path
	tmpPath := filepath.Join(uploadPath, fileName)

	// Check if the temporary file exists
	_, err := os.Stat(tmpPath)
	if err != nil {
		return errors.New("file_not_found")
	}

	// Ensure the directory exists
	fsPath, err := storage.RelativeToFsPath(relativePath)
	if err != nil {
		return err
	}

	// Make sure the directory exists
	dir := filepath.Dir(fsPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	// The target file path
	filePath := filepath.Join(fsPath, model.FileAppend)

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
	err = os.Rename(tmpPath, filePath) // This will replace the file if it already exists
	if err != nil {
		return err
	}

	// Update the file meta
	err = storage.UpdateFileMeta(fsPath, model.FileMeta{
		RelativePath: relativePath,
		FileName:     filepath.Base(relativePath),
		UploadedAt:   time.Now().Unix(),
		MimeType:     mimeTypeStr,
	})
	if err != nil {
		return err
	}

	// Update the file hash
	go hasher.HashFileAsync(fsPath)
	// Remove image cache ending with .image_*
	go image.RemoveImageCache(fsPath)

	return nil
}
