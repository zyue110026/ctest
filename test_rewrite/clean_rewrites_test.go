package testrewrite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCleanRewrites rolls back all rewritten files to original state
func TestCleanRewrites(t *testing.T) {
	k8sRoot := os.Getenv("K8S_ROOT")
	fmt.Println("K8S_ROOT:", k8sRoot)
	if k8sRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working dir: %v", err)
		}
		k8sRoot = filepath.Clean(filepath.Join(cwd, "../../.."))
	}

	// Default folder
	target := os.Getenv("CLEAN_TARGET")
	if target == "" {
		target = "test/e2e" // default folder
	}

	var absTarget string
	if filepath.IsAbs(target) {
		absTarget = target
	} else {
		absTarget = filepath.Join(k8sRoot, target)
	}

	//absTarget, err := filepath.Abs(target)
	// if err != nil {
	// 	t.Fatalf("failed to resolve absolute path: %v", err)
	// }

	files, err := collectAllGoFiles(absTarget)
	if err != nil {
		t.Fatalf("failed to collect Go files: %v", err)
	}

	deleted := 0
	cleaned := 0
	skipped := 0

	for _, f := range files {
		base := filepath.Base(f)

		// Delete rewritten files
		if strings.HasPrefix(base, "ctest_") {
			if err := os.Remove(f); err != nil {
				t.Errorf("failed to delete %s: %v", f, err)
			} else {
				t.Logf("🗑️  Deleted rewritten file: %s", f)
				deleted++
			}
			continue
		}

		// Remove build tags from original files
		hasTag, err := hasOriginalBuildTag(f)
		if err != nil {
			t.Errorf("failed checking build tag for %s: %v", f, err)
			skipped++
			continue
		}

		if hasTag {
			if err := removeOriginalBuildTag(f); err != nil {
				t.Errorf("failed removing build tag for %s: %v", f, err)
				skipped++
			} else {
				t.Logf("🏷️  Removed build tag from: %s", f)
				cleaned++
			}
		} else {
			skipped++
		}
	}

	t.Log("===================================")
	t.Logf("Clean Summary")
	t.Logf("Deleted rewritten files : %d", deleted)
	t.Logf("Cleaned build tags      : %d", cleaned)
	t.Logf("Skipped files           : %d", skipped)
	t.Log("===================================")
}

// Collect all Go files recursively, skipping certain dirs/files
func collectAllGoFiles(root string) ([]string, error) {
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

		if skipFiles[d.Name()] {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
