package task

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/util"
	"github.com/vvbbnn00/goflet/util/log"
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
				log.Infof("Remove outdated file: %s", path)
				_ = os.Remove(path)
			}
		}
		return nil
	})
}
