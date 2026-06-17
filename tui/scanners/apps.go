package scanners

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// AppSuggestion represents an application recommended to uninstall
type AppSuggestion struct {
	Name      string
	Path      string
	Size      int64
	LastUsed  time.Time
	NeverUsed bool
}

// ScanApps suggestions scans installed apps and checks last usage
func ScanApps() []AppSuggestion {
	var list []AppSuggestion
	dirs := []string{"/Applications", filepath.Join(os.Getenv("HOME"), "Applications")}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".app") {
				appPath := filepath.Join(dir, entry.Name())
				nameLower := strings.ToLower(entry.Name())

				if strings.Contains(nameLower, "safari") || strings.Contains(nameLower, "system settings") || strings.Contains(nameLower, "finder") {
					continue
				}

				size := DirSize(appPath)

				cmd := exec.Command("mdls", "-name", "kMDItemLastUsedDate", "-raw", appPath)
				var out bytes.Buffer
				cmd.Stdout = &out
				_ = cmd.Run()

				rawDate := strings.TrimSpace(out.String())
				var lastUsed time.Time
				neverUsed := false

				if rawDate == "" || rawDate == "(null)" || strings.Contains(rawDate, "could not find") {
					neverUsed = true
				} else {
					parsed, err := time.Parse("2006-01-02 15:04:05 -0700", rawDate)
					if err == nil {
						lastUsed = parsed
					} else {
						parts := strings.Split(rawDate, " ")
						if len(parts) > 0 {
							parsedFallback, err := time.Parse("2006-01-02", parts[0])
							if err == nil {
								lastUsed = parsedFallback
							} else {
								neverUsed = true
							}
						} else {
							neverUsed = true
						}
					}
				}

				isUnused := neverUsed || time.Since(lastUsed) > 30*24*time.Hour
				if size > 50*1024*1024 && isUnused {
					list = append(list, AppSuggestion{
						Name:      strings.TrimSuffix(entry.Name(), ".app"),
						Path:      appPath,
						Size:      size,
						LastUsed:  lastUsed,
						NeverUsed: neverUsed,
					})
				}
			}
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Size > list[j].Size
	})

	return list
}
