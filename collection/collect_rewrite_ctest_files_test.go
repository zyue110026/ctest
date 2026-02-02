package collection

import (
	ctestutils "k8s.io/kubernetes/test/ctest/utils"
	"os"
	"path/filepath"
	"testing"
)

// TestCollectCtestFiles collects all ctest_ Go files and copies them
// to a test folder while preserving the directory structure
func TestCollectCtestFiles(t *testing.T) {
	k8sRoot := os.Getenv("K8S_ROOT")
	t.Log("K8S_ROOT:", k8sRoot)
	if k8sRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working dir: %v", err)
		}
		k8sRoot = filepath.Clean(filepath.Join(cwd, "../../.."))
	}

	target := os.Getenv("COLLECT_TARGET")
	if target == "" {
		target = "test/integration" // default folder
	}

	var absTarget string
	if filepath.IsAbs(target) {
		absTarget = target
	} else {
		absTarget = filepath.Join(k8sRoot, target)
	}

	// Destination folder for test collection
	destRoot := `D:\k8s test rewrite collection`

	files, err := ctestutils.CollectAllGoFiles(absTarget)
	if err != nil {
		t.Fatalf("failed to collect Go files: %v", err)
	}

	copied := 0
	skipped := 0

	for _, f := range files {
		base := filepath.Base(f)
		if !startsWithCtest(base) {
			skipped++
			continue
		}

		relPath, err := filepath.Rel(k8sRoot, f)
		if err != nil {
			t.Errorf("failed to get relative path for %s: %v", f, err)
			skipped++
			continue
		}

		destPath := filepath.Join(destRoot, relPath)
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			t.Errorf("failed to create dir %s: %v", destDir, err)
			skipped++
			continue
		}

		data, err := os.ReadFile(f)
		if err != nil {
			t.Errorf("failed to read file %s: %v", f, err)
			skipped++
			continue
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			t.Errorf("failed to write file %s: %v", destPath, err)
			skipped++
			continue
		}

		t.Logf("📄 Copied: %s -> %s", f, destPath)
		copied++
	}

	t.Log("===================================")
	t.Log("Collection Summary")
	t.Logf("Copied ctest_ files : %d", copied)
	t.Logf("Skipped files       : %d", skipped)
	t.Logf("Destination root    : %s", destRoot)
	t.Log("===================================")
}

func startsWithCtest(name string) bool {
	return len(name) >= 6 && name[:6] == "ctest_"
}
