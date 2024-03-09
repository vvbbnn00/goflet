package hash

import (
	"github.com/vvbbnn00/goflet/storage/model"
	"os"
)

// getFs returns the file stream
func getFs(path string) (*os.File, error) {
	fs, err := os.OpenFile(path, os.O_RDONLY, model.FilePerm)
	if err != nil {
		return nil, err
	}
	return fs, nil
}
