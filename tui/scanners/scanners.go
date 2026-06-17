package scanners

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"syscall"
)

// DirSize recursively calculates directory size (considering sparse files via allocated blocks)
func DirSize(path string) int64 {
	var size int64
	_ = filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				if stat, ok := info.Sys().(*syscall.Stat_t); ok {
					size += stat.Blocks * 512
				} else {
					size += info.Size()
				}
			}
		}
		return nil
	})
	return size
}

// FormatSize formats bytes into human-readable MB/GB
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGT"[exp])
}
