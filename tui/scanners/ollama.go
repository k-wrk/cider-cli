package scanners

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// OllamaItem represents an Ollama model
type OllamaItem struct {
	Name string
	Path string
	Size int64
}

// ScanOllama scans ~/.ollama/models/manifests for installed models
func ScanOllama() []OllamaItem {
	var list []OllamaItem
	home, _ := os.UserHomeDir()
	manifestsDir := filepath.Join(home, ".ollama/models/manifests/registry.ollama.ai")

	_ = filepath.WalkDir(manifestsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			rel, err := filepath.Rel(manifestsDir, path)
			if err != nil {
				return nil
			}

			parts := strings.Split(rel, string(filepath.Separator))
			name := rel
			if len(parts) >= 2 {
				modelName := parts[len(parts)-2]
				tagName := parts[len(parts)-1]
				if len(parts) == 3 && parts[0] == "library" {
					name = fmt.Sprintf("%s:%s", modelName, tagName)
				} else {
					name = fmt.Sprintf("%s:%s", strings.Join(parts[:len(parts)-1], "/"), tagName)
				}
			}

			data, err := os.ReadFile(path)
			if err == nil {
				var manifest struct {
					Config struct {
						Size int64 `json:"size"`
					} `json:"config"`
					Layers []struct {
						Size int64 `json:"size"`
					} `json:"layers"`
				}
				if json.Unmarshal(data, &manifest) == nil {
					var size int64 = manifest.Config.Size
					for _, layer := range manifest.Layers {
						size += layer.Size
					}
					list = append(list, OllamaItem{
						Name: name,
						Path: path,
						Size: size,
					})
				}
			}
		}
		return nil
	})

	sort.Slice(list, func(i, j int) bool {
		return list[i].Size > list[j].Size
	})

	return list
}

// DeleteOllamaModel moves the model manifest and referenced blobs to the Trash
func DeleteOllamaModel(manifestPath string) error {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	var manifest struct {
		Config struct {
			Digest string `json:"digest"`
		} `json:"config"`
		Layers []struct {
			Digest string `json:"digest"`
		} `json:"layers"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	home, _ := os.UserHomeDir()

	moveDigestToTrash := func(digest string) {
		if digest == "" {
			return
		}
		blobName := strings.Replace(digest, ":", "-", 1)
		blobPath := filepath.Join(home, ".ollama/models/blobs", blobName)
		trashBlobPath := filepath.Join(home, ".Trash", blobName)

		if _, err := os.Stat(blobPath); err == nil {
			if _, err := os.Stat(trashBlobPath); err == nil {
				_ = os.RemoveAll(trashBlobPath)
			}
			_ = os.Rename(blobPath, trashBlobPath)
		}
	}

	for _, layer := range manifest.Layers {
		moveDigestToTrash(layer.Digest)
	}
	moveDigestToTrash(manifest.Config.Digest)

	trashManifestPath := filepath.Join(home, ".Trash", filepath.Base(manifestPath))
	if _, err := os.Stat(trashManifestPath); err == nil {
		_ = os.RemoveAll(trashManifestPath)
	}
	_ = os.Rename(manifestPath, trashManifestPath)

	return nil
}
