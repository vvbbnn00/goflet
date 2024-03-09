// Package main generates winres.json
package main

import (
	"bytes"
	"fmt"
	"github.com/vvbbnn00/goflet/base"
	"os"
	"os/exec"
	"strings"
	"time"
)

const js = `{
  "RT_GROUP_ICON": {
    "APP": {
      "0000": [
        "icon.png",
        "icon16.png"
      ]
    }
  },
  "RT_MANIFEST": {
    "#1": {
      "0409": {
        "identity": {
          "name": "go-cqhttp",
          "version": "%s"
        },
        "description": "",
        "minimum-os": "vista",
        "execution-level": "as invoker",
        "ui-access": false,
        "auto-elevate": false,
        "dpi-awareness": "system",
        "disable-theming": false,
        "disable-window-filtering": false,
        "high-resolution-scrolling-aware": false,
        "ultra-high-resolution-scrolling-aware": false,
        "long-path-aware": false,
        "printer-driver-isolation": false,
        "gdi-scaling": false,
        "segment-heap": false,
        "use-common-controls-v6": false
      }
    }
  },
  "RT_VERSION": {
    "#1": {
      "0000": {
        "fixed": {
          "file_version": "%s",
          "product_version": "%s",
          "timestamp": "%s"
        },
        "info": {
          "0409": {
            "Comments": "Goflet is a lightweight file upload and download service written in Go.",
            "CompanyName": "vvbbnn00",
            "FileDescription": "https://github.com/vvbbnn00/goflet",
            "FileVersion": "%s",
            "InternalName": "",
            "LegalCopyright": "©️ 2024 - %d vvbbnn00. All Rights Reserved.",
            "LegalTrademarks": "",
            "OriginalFilename": "GOFLET.EXE",
            "PrivateBuild": "",
            "ProductName": "Goflet",
            "ProductVersion": "%s",
            "SpecialBuild": ""
          }
        }
      }
    }
  }
}`

const timeformat = `2006-01-02T15:04:05+08:00`

func main() {
	f, err := os.Create("winres.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	v := ""
	if base.Version == "(devel)" {
		vartag := bytes.NewBuffer(nil)
		vartagcmd := exec.Command("git", "tag", "--sort=committerdate")
		vartagcmd.Stdout = vartag
		err = vartagcmd.Run()
		if err != nil {
			panic(err)
		}
		s := strings.Split(vartag.String(), "\n")
		v = s[len(s)-2]
	} else {
		v = base.Version
	}
	i := strings.Index(v, "-") // remove -rc / -beta
	if i <= 0 {
		i = len(v)
	}
	fmt.Printf("Version: %s\n", v)
	commitcnt := strings.Builder{}
	commitcnt.WriteString(v[1:i])
	commitcnt.WriteByte('.')
	commitcntcmd := exec.Command("git", "rev-list", "--count", "master")
	commitcntcmd.Stdout = &commitcnt
	err = commitcntcmd.Run()
	if err != nil {
		panic(err)
	}
	fv := commitcnt.String()[:commitcnt.Len()-1]
	_, err = fmt.Fprintf(f, js, fv, fv, v, time.Now().Format(timeformat), fv, time.Now().Year(), v)
	if err != nil {
		panic(err)
	}
}
