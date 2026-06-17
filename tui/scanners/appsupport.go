package scanners

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// AppSupportItem represents an Application Support subdirectory
type AppSupportItem struct {
	Name     string
	Path     string
	Size     int64
	Orphan   bool
	Selected bool
}

// ScanAppSupport scans Library/Application Support subdirectories larger than 10MB
func ScanAppSupport() []AppSupportItem {
	var list []AppSupportItem
	home, _ := os.UserHomeDir()
	appSupportDir := filepath.Join(home, "Library/Application Support")

	systemFolders := map[string]bool{
		"Apple": true, "com.apple.TCC": true, "SyncServices": true, "AddressBook": true,
		"Quick Look": true, "AppStore": true, "MobileSync": true, "iCloud": true,
		"com.apple.sharedfilelist": true, "com.apple.spotlight": true, "com.apple.touristd": true,
		"NotificationCenter": true, "Dock": true, "OpenDirectory": true, "DiskImages": true,
		"CrashReporter": true, "Accounts": true, "CloudDocs": true, "CoreParsec": true,
		"SearchParty": true, "SpamSieve": true, "Adobe": true, "Microsoft": true, "Oracle": true,
		"Helper": true, "com.apple.spotlight-server": true,
	}

	entries, err := os.ReadDir(appSupportDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				if systemFolders[name] || strings.HasPrefix(name, "com.apple.") {
					continue
				}

				path := filepath.Join(appSupportDir, name)
				size := DirSize(path)
				if size > 10*1024*1024 { // Only show dirs larger than 10MB
					installed := IsAppInstalled(name)
					list = append(list, AppSupportItem{
						Name:   name,
						Path:   path,
						Size:   size,
						Orphan: !installed,
					})
				}
			}
		}
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Orphan && !list[j].Orphan {
			return true
		}
		if !list[i].Orphan && list[j].Orphan {
			return false
		}
		return list[i].Size > list[j].Size
	})

	return list
}

func IsAppInstalled(folderName string) bool {
	cleanName := strings.ToLower(folderName)
	if strings.Contains(cleanName, "zoom") {
		cleanName = "zoom"
	}
	if strings.Contains(cleanName, "stremio") {
		cleanName = "stremio"
	}
	if strings.Contains(cleanName, "postman") {
		cleanName = "postman"
	}
	if strings.Contains(cleanName, "steam") {
		cleanName = "steam"
	}
	if strings.Contains(cleanName, "discord") {
		cleanName = "discord"
	}
	if strings.Contains(cleanName, "slack") {
		cleanName = "slack"
	}

	dirs := []string{"/Applications", filepath.Join(os.Getenv("HOME"), "Applications")}
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			appName := strings.ToLower(entry.Name())
			if strings.Contains(appName, cleanName) && strings.HasSuffix(appName, ".app") {
				return true
			}
		}
	}
	return false
}
