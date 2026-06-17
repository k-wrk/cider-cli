package scanners

import (
	"os"
	"path/filepath"
	"sort"
)

// BrowserItem represents a browser cache directory
type BrowserItem struct {
	Name     string
	Path     string
	Size     int64
	Selected bool
}

// ScanBrowsers scans browser cache directories
func ScanBrowsers() []BrowserItem {
	var list []BrowserItem
	home, _ := os.UserHomeDir()

	browserCaches := map[string]string{
		"Google Chrome Cache":                  filepath.Join(home, "Library/Caches/Google/Chrome"),
		"Google Chrome Support Cache":          filepath.Join(home, "Library/Application Support/Google/Chrome/Default/Cache"),
		"Brave Browser Cache":                  filepath.Join(home, "Library/Caches/BraveSoftware/Brave-Browser"),
		"Brave Browser Support Cache":          filepath.Join(home, "Library/Application Support/BraveSoftware/Brave-Browser/Default/Cache"),
		"Safari Cache":                         filepath.Join(home, "Library/Caches/com.apple.Safari"),
		"Mozilla Firefox Cache":                filepath.Join(home, "Library/Caches/Firefox"),
		"Mozilla General Cache":                filepath.Join(home, "Library/Caches/Mozilla"),
		"Zen Browser Cache":                    filepath.Join(home, "Library/Caches/Zen"),
		"DuckDuckGo Browser Cache":             filepath.Join(home, "Library/Caches/com.duckduckgo.macos.browser"),
	}

	for name, path := range browserCaches {
		if _, err := os.Stat(path); err == nil {
			size := DirSize(path)
			if size > 0 {
				list = append(list, BrowserItem{
					Name: name,
					Path: path,
					Size: size,
				})
			}
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Size > list[j].Size
	})

	return list
}
