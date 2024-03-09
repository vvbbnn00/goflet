package image

import (
	"bytes"
	"goflet/storage/model"
	"io"
	"log"
	"os"
	"path/filepath"
)

// GetFileImageReader get the file reader for the image
func GetFileImageReader(fsPath string, params *ProcessParams) (*os.File, error) {
	fsPath = filepath.Join(fsPath, model.ImageAppend+params.Dump())

	file, err := os.OpenFile(fsPath, os.O_RDONLY, model.FilePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SaveFileImageCache save the file to the image cache
func SaveFileImageCache(fsPath string, params *ProcessParams, buffer bytes.Buffer) error {
	fsPath = filepath.Join(fsPath, model.ImageAppend+params.Dump())

	// Copy the file to the cache
	cacheFile, err := os.OpenFile(fsPath, os.O_CREATE|os.O_RDWR, model.FilePerm)
	if err != nil {
		_ = cacheFile.Close()
		return err
	}

	// Write the buffer to the file
	_, err = io.Copy(cacheFile, &buffer)
	if err != nil {
		_ = cacheFile.Close()
		return err
	}

	// Close the file
	_ = cacheFile.Close()

	return nil
}

// RemoveImageCache remove the image cache
func RemoveImageCache(fsPath string) {
	// Remove the file from the cache
	pathPattern := filepath.Join(fsPath, model.ImageCachePrefixWithWildcard)
	files, err := filepath.Glob(pathPattern)

	if err != nil {
		log.Printf("Error removing image cache: %s", err.Error())
	}

	// Remove the files
	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			continue
		}
	}
}
