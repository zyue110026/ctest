package fixtures

import (
	"embed"
	"encoding/json"
	"fmt"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	"strings"
)

//go:embed *.json
var fixtureFS embed.FS

// LoadFixtures loads JSON from embedded fixtures and returns only the requested types.
// If no types are provided, it returns all non-null top-level keys.
// Returned map maps key -> json.RawMessage (the JSON subtree for that key).
//
// Behavior:
// - If requested types are passed and any of them are not present in the file -> returns an error listing missing keys.
// - If a requested key exists but its value is JSON null, it is omitted from the returned map (no error).
// - If no types requested: include all top-level keys whose value != null.
func LoadFixturesAsJSON(fileName string, types ...string) (map[string]json.RawMessage, error) {
	// Construct the path within the embedded filesystem
	// fsPath := "fixtures/" + fileName
	fmt.Println(ctestglobals.DebugPrefix(), "Loading embedded fixture file:", fileName)

	// Read file from embedded filesystem
	b, err := fixtureFS.ReadFile(fileName)
	if err != nil {
		fmt.Println(ctestglobals.DebugPrefix(), "Error reading embedded file:", err)
		return nil, fmt.Errorf("read embedded file: %w", err)
	}

	// Unmarshal into a map of raw messages so we can inspect each top-level value
	var root map[string]json.RawMessage
	if err := json.Unmarshal(b, &root); err != nil {
		fmt.Println(ctestglobals.DebugPrefix(), "Error unmarshaling json:", err)
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	result := make(map[string]json.RawMessage)

	// helper to check if raw json is null
	isNull := func(r json.RawMessage) bool {
		trim := strings.TrimSpace(string(r))
		return trim == "" || trim == "null"
	}

	// fmt.Printf(ctestglobals.DebugPrefix(), "Requested fixture types: %v\n", types)
	// fmt.Println(ctestglobals.DebugPrefix(), "Has types requested:", len(types) > 0)

	if len(types) == 0 {
		// return all keys whose value is not null
		for k, v := range root {
			if !isNull(v) {
				result[k] = v
			}
		}
		return result, nil
	}

	// When specific types requested, ensure they exist in the file.
	var missing []string
	for _, t := range types {
		v, ok := root[t]
		if !ok {
			missing = append(missing, t)
			continue
		}
		// If it exists but is null, skip (no error)
		if !isNull(v) {
			result[t] = v
		}
	}

	if len(missing) > 0 {
		fmt.Println(ctestglobals.DebugPrefix(), "Missing requested fixture keys:", missing)
		return nil, fmt.Errorf("requested fixture keys not found in %s: %s", fileName, strings.Join(missing, ", "))
	}
	return result, nil
}

// Optional: Add a function to list all available fixture files
func ListFixtureFiles() ([]string, error) {
	entries, err := fixtureFS.ReadDir("fixtures")
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}
