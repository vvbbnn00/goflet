package upload

import (
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"goflet/cache"
	"goflet/config"
	"goflet/storage"
	"goflet/storage/hasher"
	"goflet/storage/image"
	"goflet/storage/model"
	"goflet/util"
	"goflet/util/hash"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	cachePrefix = "uploading:"
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
	c := cache.GetCache()
	// Ensure the directory exists
	fsPath, err := storage.RelativeToFsPath(relativePath)
	if err != nil {
		return err
	}

	exists, _ := c.GetBool(cachePrefix + fsPath)
	if exists {
		return errors.New("file_uploading")
	}

	// Check if the temporary file exists
	_, err = os.Stat(tmpPath)
	if err != nil {
		return errors.New("file_not_found")
	}

	// Make sure the directory exists
	dir := filepath.Dir(fsPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

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

	// Complete the upload
	go completeUpload(fsPath, tmpPath, model.FileMeta{
		RelativePath: relativePath,
		FileName:     filepath.Base(relativePath),
		MimeType:     mimeTypeStr,
		UploadedAt:   time.Now().Unix(),
	})

	return nil
}

// completeUpload completes the file upload by renaming the temporary file to the final file
func completeUpload(fsPath string, tmpPath string, meta model.FileMeta) {
	c := cache.GetCache()
	_ = c.SetEx(cachePrefix+fsPath, true, 60)

	defer func() {
		_ = c.Del(cachePrefix + fsPath)
	}()

	// The target file path
	filePath := filepath.Join(fsPath, model.FileAppend)

	// Rename the temporary file to the final file
	err := storage.MoveFile(tmpPath, filePath)
	if err != nil {
		log.Printf("Error moving file: %s", err.Error())
		return // Give up if the file cannot be moved
	}

	// Update the file meta
	err = storage.UpdateFileMeta(fsPath, meta)
	if err != nil {
		log.Printf("Error updating file meta: %s", err.Error())
		return // Give up if the file meta cannot be updated
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Update the file hash
	go func() {
		hasher.HashFileAsync(fsPath)
		wg.Done()
	}()
	// Remove image cache ending with .image_*
	go func() {
		image.RemoveImageCache(fsPath)
		wg.Done()
	}()

	wg.Wait()
}
