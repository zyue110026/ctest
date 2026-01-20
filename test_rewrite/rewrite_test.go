package testrewrite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestRewriteWithLLM rewrites Go test files using Ollama (DeepSeek-Coder)
func TestRewriteWithLLM(t *testing.T) {
	start := time.Now()

	k8sRoot := os.Getenv("K8S_ROOT")
	fmt.Println("K8S_ROOT:", k8sRoot)
	if k8sRoot == "" {
		var err error
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working dir: %v", err)
		}
		// go two levels up from current working dir
		k8sRoot = filepath.Clean(filepath.Join(cwd, "../.."))
	}

	target := os.Getenv("REWRITE_TARGET")
	if target == "" {
		target = "test/e2e" // default
	}

	// resolve absolute path
	var absTarget string
	if filepath.IsAbs(target) {
		absTarget = target
	} else {
		absTarget = filepath.Join(k8sRoot, target)
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "deepseek-coder:33b"
	}

	t.Logf("Rewrite target: %s", absTarget)
	t.Logf("Using Ollama model: %s", model)

	files, err := CollectGoFiles(absTarget)
	if err != nil {
		t.Fatalf("failed to collect files: %v", err)
	}

	if len(files) == 0 {
		t.Log("No Go files to rewrite")
		return
	}

	for i, file := range files {
		t.Logf("[%d/%d] Rewriting %s", i+1, len(files), file)

		contentBytes, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("failed to read file %s: %v", file, err)
			continue
		}

		content := string(contentBytes)
		rewritten, err := CallOllamaWithChunks(file, content, 150)
		if err != nil {
			t.Errorf("rewrite failed for %s: %v", file, err)
			continue
		}

		if strings.EqualFold(strings.TrimSpace(rewritten), "NONE") {
			t.Logf("⚠️ No tests need rewriting in %s, skipping", file)
			continue
		}

		newFile := filepath.Join(filepath.Dir(file), "ctest_"+filepath.Base(file))
		if err := os.WriteFile(newFile, []byte(rewritten), 0644); err != nil {
			t.Errorf("failed to write rewritten file %s: %v", newFile, err)
			continue
		}

		t.Logf("✅ Rewritten file saved as %s", newFile)
	}

	t.Logf("Rewrite finished in %s", time.Since(start))
}
