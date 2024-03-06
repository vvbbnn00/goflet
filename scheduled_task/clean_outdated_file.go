package scheduled_task

import (
	"goflet/config"
	"goflet/util"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
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
		if len(entries) == 0 {
			log.Printf("Remove empty folder: %s", path)
			_ = os.Remove(path)
		}
	}
}

// CleanOutdatedFile Clean outdated files
func CleanOutdatedFile() {
	uploadPath := util.GetUploadPath()
	UploadTimeout := time.Duration(config.GofletCfg.FileConfig.UploadTimeout) * time.Second

	_ = filepath.Walk(uploadPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if time.Since(info.ModTime()) > UploadTimeout {
				log.Printf("Remove outdated file: %s", path)
				_ = os.Remove(path)
			}
		}
		return nil
	})
}
