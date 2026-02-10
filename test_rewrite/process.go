package testrewrite

import (
	"fmt"

	"os"
	"path/filepath"
	"strings"
)

var repoRootAllowDirs = map[string]bool{
	"api":     false,
	"cluster": true,
	"cmd":     true,
	"pkg":     true,
	"plugin":  true,
	"staging": true,
	"test":    true, // must keep this
}

var skipDirs = map[string]bool{
	"framework": true,
	"utils":     true,
	"ctest":     true,
	".github":   true,
}

var skipFiles = map[string]bool{
	"framework.go": true,
	"utils.go":     true,
	"doc.go":       true,
	"util.go":      true,
}

func CollectGoFilesFromRepo(k8sRoot, path string) ([]string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	testRoot := filepath.Join(k8sRoot, "test")

	// ✅ Inside k8sRoot/test/** → collect all go files with skip rules
	if isUnder(absPath, testRoot) {
		return CollectGoFiles(absPath)
	}

	// ❌ Outside test/** → collect only *_test.go files with skip rules, and only if they are under allowed top-level dirs
	return walkOutsideTest(k8sRoot, absPath)
}

func isUnder(path, parent string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && !strings.HasPrefix(rel, "..")
}

func CollectGoFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}

	// Directory → walk normally
	if info.IsDir() {
		return walkTestDir(path)
	}

	// Single Go file → apply all skip rules
	if strings.HasSuffix(path, ".go") {
		base := filepath.Base(path)

		// Skip specific filenames
		if skipFiles[base] {
			return []string{}, nil
		}

		// Skip rewritten files
		if strings.HasPrefix(base, "ctest_") {
			return []string{}, nil
		}

		// Skip files inside ctest directory
		if strings.Contains(path, string(filepath.Separator)+"ctest"+string(filepath.Separator)) {
			return []string{}, nil
		}

		return []string{path}, nil
	}

	return nil, fmt.Errorf("unsupported target: %s", path)
}

func walkTestDir(root string) ([]string, error) {
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
		// fmt.Println(d.Name())
		// Skip framework/utils and already rewritten files
		if skipFiles[d.Name()] ||
			strings.HasPrefix(d.Name(), "ctest_") ||
			strings.Contains(d.Name(), "ctest_") {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

func walkOutsideTest(k8sRoot, root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(k8sRoot, path)
		if err != nil {
			return err
		}

		parts := strings.Split(rel, string(filepath.Separator))

		// --------------------
		// Directory rules
		// --------------------
		if d.IsDir() {
			// Repo-root level: ONLY allow specific dirs
			if len(parts) == 1 && !repoRootAllowDirs[parts[0]] {
				fmt.Printf(
					"    process.go:145: [collect][repo-root-skip] skipping directory: %s\n",
					parts[0],
				)
				return filepath.SkipDir
			}
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		name := d.Name()

		// --------------------
		// File rules
		// --------------------

		// Only *_test.go
		if !strings.HasSuffix(name, "_test.go") {
			return nil
		}

		// Skip rewritten files
		if strings.HasPrefix(name, "ctest_") {
			return nil
		}

		// Reuse existing skip rules
		if skipFiles[name] {
			return nil
		}

		if strings.Contains(path, string(filepath.Separator)+"ctest"+string(filepath.Separator)) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
