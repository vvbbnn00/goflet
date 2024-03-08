package image

import (
	"bytes"
	"goflet/service"
	"io"
	"os"
)

const (
	FileAppend = ".image_"
)

// GetFileImageReader get the file reader for the image
func GetFileImageReader(path string, params *ProcessParams) (*os.File, error) {
	fsPath, err := service.ConvertToFsPath(path)
	if err != nil {
		return nil, err
	}

	fsPath += FileAppend + params.Dump()

	file, err := os.OpenFile(fsPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SaveFileImageCache save the file to the image cache
func SaveFileImageCache(path string, params *ProcessParams, buffer bytes.Buffer) error {
	fsPath, err := service.ConvertToFsPath(path)
	if err != nil {
		return err
	}

	fsPath += FileAppend + params.Dump()

	// Copy the file to the cache
	cacheFile, err := os.OpenFile(fsPath, os.O_CREATE|os.O_RDWR, 0644)
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

	println("File saved to cache: " + fsPath)

	return nil
}
