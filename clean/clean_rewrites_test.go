package clean

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ctestutils "k8s.io/kubernetes/test/ctest/utils"
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

	files, err := ctestutils.CollectAllGoFiles(absTarget)
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

func hasOriginalBuildTag(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	content := string(data)

	return strings.HasPrefix(content, "//go:build original") ||
		strings.HasPrefix(content, "// +build original"), nil
}

func removeOriginalBuildTag(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "//go:build original") || strings.HasPrefix(line, "// +build original") {
			continue
		}
		newLines = append(newLines, line)
	}

	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
}
