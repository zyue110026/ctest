package ctest

import (
	"io/fs"
	//"os"
	"path/filepath"
	"strings"

	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
)

func shouldSkipPath(path string) bool {
	lower := strings.ToLower(path)

	// skip templates directory
	if strings.Contains(lower, string(filepath.Separator)+"templates"+string(filepath.Separator)) {
		return true
	}

	for _, w := range ctestglobals.WeirdPaths {
		if strings.Contains(lower, strings.ToLower(w)) {
			return true
		}
	}
	return false
}

func collectYAMLFiles(repo string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(repo, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // soft-fail
		}

		if shouldSkipPath(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
