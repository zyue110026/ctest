package testrewrite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var skipDirs = map[string]bool{
	"framework": true,
	"utils":     true,
	"ctest":     true,
}

var skipFiles = map[string]bool{
	"framework.go": true,
	"utils.go":     true,
}

func CollectGoFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}

	if info.IsDir() {
		return walkDir(path)
	}

	if strings.HasSuffix(path, ".go") {
		return []string{path}, nil
	}

	return nil, fmt.Errorf("unsupported target: %s", path)
}

func walkDir(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		// Skip framework/utils and already rewritten files
		if skipFiles[d.Name()] || strings.HasPrefix(d.Name(), "ctest_") {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
