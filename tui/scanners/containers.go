package scanners

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// List of protected whitelist keywords to prevent false-positives in containers
var safeKeywords = []string{
	"apple", "microsoft", "office", "onedrive", "adobe", "docker", "1password",
	"dropbox", "google", "chrome", "spotify", "slack", "telegram", "discord",
	"visualstudio", "vscode", "cursor", "intellij", "github", "brave", "safari",
}

// IsSafeToClean analyzes if a container folder is in the whitelist
func IsSafeToClean(folderName string) bool {
	name := strings.ToLower(folderName)
	for _, word := range safeKeywords {
		if strings.Contains(name, word) {
			return false
		}
	}
	return true
}

// GetOrphanContainersSize calculates the size of all suspected orphan containers
func GetOrphanContainersSize() int64 {
	home, _ := os.UserHomeDir()
	dirs := []string{
		filepath.Join(home, "Library/Containers"),
		filepath.Join(home, "Library/Group Containers"),
	}

	var totalSize int64
	for _, baseDir := range dirs {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				folderName := entry.Name()
				if !IsSafeToClean(folderName) {
					continue
				}
				cleanName := folderName
				if strings.HasPrefix(cleanName, "group.") {
					cleanName = strings.TrimPrefix(cleanName, "group.")
				}
				parts := strings.SplitN(cleanName, ".", 2)
				if len(parts) == 2 && len(parts[0]) == 10 {
					cleanName = parts[1]
				}

				cmd1 := exec.Command("mdfind", fmt.Sprintf("kMDItemCFBundleIdentifier == %s", folderName))
				var out1 bytes.Buffer
				cmd1.Stdout = &out1
				_ = cmd1.Run()

				cmd2 := exec.Command("mdfind", fmt.Sprintf("kMDItemCFBundleIdentifier == %s", cleanName))
				var out2 bytes.Buffer
				cmd2.Stdout = &out2
				_ = cmd2.Run()

				if !strings.Contains(out1.String(), ".app") && !strings.Contains(out2.String(), ".app") {
					totalSize += DirSize(filepath.Join(baseDir, folderName))
				}
			}
		}
	}
	return totalSize
}

// CleanOrphanContainers moves orphan containers to macOS Trash
func CleanOrphanContainers() (string, error) {
	home, _ := os.UserHomeDir()
	dirs := []string{
		filepath.Join(home, "Library/Containers"),
		filepath.Join(home, "Library/Group Containers"),
	}

	var count int
	for _, baseDir := range dirs {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				folderName := entry.Name()
				if !IsSafeToClean(folderName) {
					continue
				}

				cleanName := folderName
				if strings.HasPrefix(cleanName, "group.") {
					cleanName = strings.TrimPrefix(cleanName, "group.")
				}
				parts := strings.SplitN(cleanName, ".", 2)
				if len(parts) == 2 && len(parts[0]) == 10 {
					cleanName = parts[1]
				}

				cmd1 := exec.Command("mdfind", fmt.Sprintf("kMDItemCFBundleIdentifier == %s", folderName))
				var out1 bytes.Buffer
				cmd1.Stdout = &out1
				_ = cmd1.Run()

				cmd2 := exec.Command("mdfind", fmt.Sprintf("kMDItemCFBundleIdentifier == %s", cleanName))
				var out2 bytes.Buffer
				cmd2.Stdout = &out2
				_ = cmd2.Run()

				if !strings.Contains(out1.String(), ".app") && !strings.Contains(out2.String(), ".app") {
					trashPath := filepath.Join(home, ".Trash", folderName)
					_ = os.Rename(filepath.Join(baseDir, folderName), trashPath)
					count++
				}
			}
		}
	}
	return fmt.Sprintf("%d orphan containers moved to the Trash", count), nil
}
