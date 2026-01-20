package ctest

import (
	"flag"
	"testing"

	"k8s.io/kubernetes/test/ctest/fixtures"
)

var repoDir string

func init() {
	flag.StringVar(&repoDir, "repo", "", "path to repo root")
}

func TestGenerateFixtures(t *testing.T) {
	flag.Parse()

	if repoDir == "" {
		t.Fatal("missing -repo flag")
	}

	// reset fixture store
	if err := fixtures.ClearFixtures(); err != nil {
		t.Fatal(err)
	}

	files, err := collectYAMLFiles(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	var allObjects []K8sObject

	for _, f := range files {
		objs, err := parseYAMLFile(f)
		if err != nil {
			continue
		}
		allObjects = append(allObjects, objs...)
	}

	if len(allObjects) == 0 {
		t.Fatal("no valid kubernetes objects found")
	}

	ProcessObjects(allObjects)
}
