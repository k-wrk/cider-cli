package scanners

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HFItem represents a Hugging Face cache item (model or other)
type HFItem struct {
	Name string
	Path string
	Size int64
}

// ScanHuggingFace scans ~/.cache/huggingface for models and cache items
func ScanHuggingFace() []HFItem {
	var list []HFItem
	home, _ := os.UserHomeDir()
	hfDir := filepath.Join(home, ".cache/huggingface")

	// Scan hub for models
	hubDir := filepath.Join(hfDir, "hub")
	if entries, err := os.ReadDir(hubDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), "models--") {
				path := filepath.Join(hubDir, entry.Name())
				size := DirSize(path)
				name := parseHFModelName(entry.Name())
				list = append(list, HFItem{
					Name: "Model: " + name,
					Path: path,
					Size: size,
				})
			}
		}
	}

	// Scan other directories under ~/.cache/huggingface/
	if entries, err := os.ReadDir(hfDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() && entry.Name() != "hub" {
				path := filepath.Join(hfDir, entry.Name())
				size := DirSize(path)
				list = append(list, HFItem{
					Name: "Cache: " + entry.Name(),
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

func parseHFModelName(folderName string) string {
	if strings.HasPrefix(folderName, "models--") {
		name := strings.TrimPrefix(folderName, "models--")
		return strings.ReplaceAll(name, "--", "/")
	}
	return folderName
}
