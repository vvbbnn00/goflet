package scheduled_task

import (
	"goflet/util"
	"goflet/util/log"
	"io/fs"
	"os"
	"path/filepath"
)

// DeleteEmptyFolder Delete empty folders
func DeleteEmptyFolder() {
	dataPath := util.GetBasePath()

	var pathToCheckList []string

	// Recursively delete empty folders
	_ = filepath.WalkDir(dataPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		pathToCheckList = append(pathToCheckList, path)
		return nil
	})

	// Check the folders in reverse order
	for i := len(pathToCheckList) - 1; i >= 0; i-- {
		path := pathToCheckList[i]
		pathFs := os.DirFS(path)
		entries, _ := fs.ReadDir(pathFs, ".")
		if len(entries) == 0 && path != dataPath {
			log.Infof("Remove empty folder: %s", path)
			_ = os.Remove(path)
		}
	}
}
