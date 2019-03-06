package baslib

import (
	"path/filepath"
)

func Files(pattern string) {
	files, errFiles := filepath.Glob(pattern)
	if errFiles != nil {
		alert("FILES %s: error: %v", pattern, errFiles)
	}
	for _, f := range files {
		Println(f)
	}
}
