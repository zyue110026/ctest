package fixtures

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
)

func TestLoadFixturesAsJSON(t *testing.T) {
	// 1) Load all non-null fixtures
	all, err := LoadFixturesAsJSON(ctestglobals.TestExternalFixtureFile)
	if err != nil {
		log.Fatalf("load all fixtures failed: %v", err)
	}
	fmt.Println("All non-null keys:", keys(all))
	//print all as pretty JSON
	pretty, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(pretty))

	// 2) Load specific types
	selected, err := LoadFixturesAsJSON(ctestglobals.TestExternalFixtureFile, "deployments", "services", "pods")
	if err != nil {
		// if "pods" not present in file this will be an error
		log.Fatalf("loading selected fixtures failed: %v", err)
	}

	fmt.Println("Selected keys:", keys(selected))
}

// helper to show keys (not required)
func keys(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
