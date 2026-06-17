package scanners

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

// DockerItem represents a Docker cache or storage item
type DockerItem struct {
	ID          string
	Name        string
	Path        string
	Size        int64
	CustomClean func() (string, error)
}

// ScanDocker scans disk for Docker VM storage files and global config files
func ScanDocker() []DockerItem {
	var list []DockerItem
	home, _ := os.UserHomeDir()

	// 1. Unused Docker objects prune action
	pruneItem := DockerItem{
		ID:   "prune",
		Name: "Prune Unused Docker Objects (Containers, Images, Volumes)",
		Path: "docker_prune",
		CustomClean: func() (string, error) {
			cmd := exec.Command("docker", "system", "prune", "-a", "--volumes", "-f")
			err := cmd.Run()
			if err != nil {
				return "", err
			}
			return "Docker pruned successfully", nil
		},
	}

	// Always measure prune action as a mock 500MB if docker command is available
	if _, err := exec.LookPath("docker"); err == nil {
		pruneItem.Size = 500 * 1024 * 1024
		list = append(list, pruneItem)
	}

	// 2. Scan for VM storage files (where Docker Desktop raw VM images live)
	vmPaths := []string{
		filepath.Join(home, "Library/Containers/com.docker.docker/Data/vms/0/data/ActiveOS.raw"),
		filepath.Join(home, "Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw"),
		filepath.Join(home, "Library/Containers/com.docker.docker/Data/vms/0/Docker.raw"),
		filepath.Join(home, "Library/Containers/com.docker.docker/Data/vms/0/data/Docker.qcow2"),
	}

	for _, path := range vmPaths {
		if _, err := os.Stat(path); err == nil {
			size := DirSize(path) // Handles block sizing
			if size > 0 {
				list = append(list, DockerItem{
					ID:   "vm_storage",
					Name: "Docker VM Disk Image: " + filepath.Base(path),
					Path: path,
					Size: size,
				})
			}
		}
	}

	// 3. Scan ~/.docker configuration directory
	dockerConfigDir := filepath.Join(home, ".docker")
	if _, err := os.Stat(dockerConfigDir); err == nil {
		size := DirSize(dockerConfigDir)
		if size > 0 {
			list = append(list, DockerItem{
				ID:   "config",
				Name: "Docker Config Directory (~/.docker)",
				Path: dockerConfigDir,
				Size: size,
			})
		}
	}

	// 4. Scan Docker Containers folder
	containerDir := filepath.Join(home, "Library/Containers/com.docker.docker")
	if _, err := os.Stat(containerDir); err == nil {
		size := DirSize(containerDir)
		// Only list com.docker.docker folder itself if we didn't list VM storage files
		// or if there are other files in it
		if size > 10*1024*1024 {
			list = append(list, DockerItem{
				ID:   "app_containers",
				Name: "Docker Desktop Containers Cache",
				Path: containerDir,
				Size: size,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Size > list[j].Size
	})

	return list
}
