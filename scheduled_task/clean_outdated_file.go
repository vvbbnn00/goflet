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
