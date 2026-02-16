package testrewrite

import (
	"fmt"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestRewriteWithLLM rewrites Go test files using Ollama
func TestRewriteWithLLM(t *testing.T) {
	start := time.Now()

	//---------------------------------------
	// Config
	//---------------------------------------

	k8sRoot := os.Getenv("K8S_ROOT")
	fmt.Println("K8S_ROOT:", k8sRoot)
	if k8sRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get current working dir: %v", err)
		}
		k8sRoot = filepath.Clean(filepath.Join(cwd, "../../.."))
	}

	target := os.Getenv("REWRITE_TARGET")
	if target == "" {
		target = "staging/src/k8s.io/sample-apiserver" //
	}

	var absTarget string
	if filepath.IsAbs(target) {
		absTarget = target
	} else {
		absTarget = filepath.Join(k8sRoot, target)
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = ctestglobals.OllamaModelDefault
	}

	overwrite := strings.EqualFold(os.Getenv("OVERWRITE_REWRITTEN"), "true")

	t.Logf("Rewrite target: %s", absTarget)
	t.Logf("Using Ollama model: %s", model)
	t.Logf("Overwrite rewritten files: %v", overwrite)

	//---------------------------------------
	// Collect files
	//---------------------------------------

	// files, err := CollectGoFiles(absTarget)
	files, err := CollectGoFilesFromRepo(k8sRoot, absTarget)

	if err != nil {
		t.Fatalf("failed to collect files: %v", err)
	}

	if len(files) == 0 {
		t.Log("No Go files to rewrite")
		return
	}

	//---------------------------------------
	// Summary counters
	//---------------------------------------

	total := len(files)
	rewritten := 0
	skipped := 0
	failed := 0
	alreadyRewritten := 0

	//---------------------------------------
	// Rewrite loop
	//---------------------------------------

	for i, file := range files {

		newFile := rewrittenPath(file)
		exists := fileExists(newFile)

		tagged, err := hasOriginalBuildTag(file)
		if err != nil {
			t.Errorf("failed checking build tag for %s: %v", file, err)
			failed++
			continue
		}

		// Skip if already rewritten and not overwriting
		if exists && !overwrite {
			t.Logf("⏭️  Skipping already rewritten file: %s", file)
			alreadyRewritten++
			continue
		}

		t.Logf("[%d/%d] Rewriting %s", i+1, total, file)

		contentBytes, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("failed to read file %s: %v", file, err)
			failed++
			continue
		}

		prompt := BuildPrompt(file, string(contentBytes))
		rewrittenContent, err := CallOllama(prompt)
		// time.Sleep(30)
		if err != nil {
			t.Errorf("rewrite failed for %s: %v", file, err)
			failed++
			continue
		}

		rewrittenTrim := strings.TrimSpace(rewrittenContent)

		// --- Handle override + NONE ---
		if strings.EqualFold(rewrittenTrim, "NONE") {
			t.Logf("⚠️  No tests need rewriting in %s", file)
			skipped++

			if overwrite && exists {
				// Delete previous rewritten file
				if err := os.Remove(newFile); err != nil {
					t.Errorf("failed to remove old rewritten file %s: %v", newFile, err)
				} else {
					t.Logf("🗑️  Deleted previous rewritten file %s", newFile)
				}

				// Remove original build tag
				if tagged {
					if err := removeOriginalBuildTag(file); err != nil {
						t.Errorf("failed to remove build tag from %s: %v", file, err)
					} else {
						t.Logf("🏷️ Removed build tag from original file %s", file)
					}
				}
			}
			continue
		}
		// --- End override + NONE ---

		// Write rewritten content
		if err := os.WriteFile(newFile, []byte(rewrittenContent), 0644); err != nil {
			t.Errorf("failed to write %s: %v", newFile, err)
			failed++
			continue
		}

		// // Add build tag to original file
		// if err := addOriginalBuildTag(file); err != nil {
		// 	t.Errorf("failed to tag original file %s: %v", file, err)
		// 	failed++
		// 	continue
		// }

		t.Logf("✅ Saved %s", newFile)
		rewritten++
	}

	//---------------------------------------
	// Final summary
	//---------------------------------------

	t.Log("===================================")
	t.Logf("Rewrite Summary")
	t.Logf("Targeted files    : %d", total)
	t.Logf("Rewritten         : %d", rewritten)
	t.Logf("Already Rewritten : %d", alreadyRewritten)
	t.Logf("Skipped           : %d", skipped)
	t.Logf("Failed            : %d", failed)
	t.Logf("Elapsed time      : %s", time.Since(start))
	t.Log("===================================")
}

//------------------------------------------------
// Helpers
//------------------------------------------------

func rewrittenPath(original string) string {
	return filepath.Join(
		filepath.Dir(original),
		"ctest_"+filepath.Base(original),
	)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

func addOriginalBuildTag(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)

	if strings.HasPrefix(content, "//go:build original") ||
		strings.HasPrefix(content, "// +build original") {
		return nil
	}

	tag := "//go:build original\n// +build original\n\n"
	return os.WriteFile(path, []byte(tag+content), 0644)
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
