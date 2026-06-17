package scanners

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileInfo represents size information of a file in the documents report
type FileInfo struct {
	Path string
	Size int64
}

// ScanDocuments scans the ~/Documents folder for the largest files
func ScanDocuments() []FileInfo {
	home, _ := os.UserHomeDir()
	docDir := filepath.Join(home, "Documents")
	var files []FileInfo

	_ = filepath.WalkDir(docDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil && !strings.HasPrefix(d.Name(), ".") {
				files = append(files, FileInfo{
					Path: path,
					Size: info.Size(),
				})
			}
		}
		return nil
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(files) > 50 {
		files = files[:50]
	}
	return files
}
