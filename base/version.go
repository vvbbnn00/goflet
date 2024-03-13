// Package base provides the base information for the application
package base

import (
	"fmt"
	"runtime/debug"
)

// Version The version of the application
var Version = "unknown"

func init() {
	if Version != "unknown" {
		return
	}
	info, ok := debug.ReadBuildInfo()
	if ok {
		Version = info.Main.Version
	}
}

// PrintBanner Print the banner
func PrintBanner() {
	fmt.Printf(`
   ___          __  _        _   
  / _ \  ___   / _|| |  ___ | |_ 
 / /_\/ / _ \ | |_ | | / _ \| __|
/ /_\\ | (_) ||  _|| ||  __/| |_ 
\____/  \___/ |_|  |_| \___| \__|

Goflet version: %s

「さぁ、始まるザマスよ！」

`, Version)
}
