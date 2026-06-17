package scanners

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

// DevItem represents a developer tools cache item
type DevItem struct {
	ID          string
	Name        string
	Path        string
	Size        int64
	CustomClean func() (string, error)
}

// ScanDevTools scans developer cache directories
func ScanDevTools() []DevItem {
	var list []DevItem
	home, _ := os.UserHomeDir()

	devPaths := []DevItem{
		{ID: "xcode", Name: "Xcode Derived Data", Path: filepath.Join(home, "Library/Developer/Xcode/DerivedData")},
		{ID: "xcode_device_support", Name: "Xcode iOS DeviceSupport Symbols", Path: filepath.Join(home, "Library/Developer/Xcode/iOS DeviceSupport")},
		{ID: "xcode_simulators", Name: "Xcode iOS Simulators Cache", Path: filepath.Join(home, "Library/Developer/CoreSimulator/Devices")},
		{ID: "xcode_archives", Name: "Xcode Archives", Path: filepath.Join(home, "Library/Developer/Xcode/Archives")},
		{ID: "go_mod", Name: "Go Module Cache", Path: filepath.Join(home, "go/pkg/mod")},
		{ID: "go_build", Name: "Go Build Cache", Path: filepath.Join(home, "Library/Caches/go-build")},
		{ID: "cargo_registry", Name: "Rust Cargo Registry Cache", Path: filepath.Join(home, ".cargo/registry/cache")},
		{ID: "cargo_git", Name: "Rust Cargo Git Cache", Path: filepath.Join(home, ".cargo/git/db")},
		{ID: "pip", Name: "Python Pip Cache", Path: filepath.Join(home, "Library/Caches/pip")},
		{ID: "poetry", Name: "Python Poetry Cache", Path: filepath.Join(home, "Library/Caches/pypoetry")},
		{ID: "cypress", Name: "Cypress Browser Binaries Cache", Path: filepath.Join(home, "Library/Caches/Cypress")},
		{ID: "jetbrains", Name: "JetBrains IDEs Caches", Path: filepath.Join(home, "Library/Caches/JetBrains")},
		{ID: "android_sdk", Name: "Android SDK Cache (AVD/Simulators)", Path: filepath.Join(home, "Library/Android")},
		{ID: "android_avds", Name: "Android Virtual Devices (AVD Files)", Path: filepath.Join(home, ".android/avd")},
		{ID: "yarn", Name: "Yarn Package Cache (Node.js)", Path: filepath.Join(home, "Library/Caches/Yarn")},
		{ID: "npm", Name: "npm Cache", Path: filepath.Join(home, ".npm/_cacache")},
		{ID: "uv_cache", Name: "Python uv Package Cache", Path: filepath.Join(home, ".cache/uv"), CustomClean: func() (string, error) {
			cmd := exec.Command("uv", "cache", "clean")
			_ = cmd.Run()
			return "uv cache cleared", nil
		}},
		{ID: "phpactor", Name: "Phpactor Caches (PHP Tooling)", Path: filepath.Join(home, ".cache/phpactor")},
		{ID: "gradle", Name: "Gradle Caches (Java/Android Build)", Path: filepath.Join(home, ".gradle/caches")},
		{ID: "cocoapods", Name: "CocoaPods Cache", Path: filepath.Join(home, "Library/Caches/CocoaPods")},
		{ID: "duzzy_osrm", Name: "Duzzy Server OSRM Data", Path: filepath.Join(home, "Project/Gokzel/duzzy-server/data")},
		{ID: "brew", Name: "Homebrew System Cleanup", Path: "brew_cleanup", CustomClean: func() (string, error) {
			cmd := exec.Command("brew", "cleanup", "-s")
			_ = cmd.Run()
			return "Homebrew cleanup completed", nil
		}},
	}

	for _, item := range devPaths {
		var size int64
		if item.ID == "brew" {
			if _, err := exec.LookPath("brew"); err == nil {
				size = 200 * 1024 * 1024
			}
		} else {
			if _, err := os.Stat(item.Path); err == nil {
				size = DirSize(item.Path)
			}
		}

		if size > 0 {
			item.Size = size
			list = append(list, item)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Size > list[j].Size
	})

	return list
}
