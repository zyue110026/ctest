package ctest

import (
	"log"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/kubernetes/test/ctest/ctestglobals"
)

func parseYAMLFile(path string) ([]K8sObject, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var objects []K8sObject
	docs := strings.Split(string(data), "---")

	for i, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		obj, gvk, err := scheme.Codecs.UniversalDeserializer().
			Decode([]byte(doc), nil, nil)

		if err != nil {
			log.Printf("skip invalid yaml %s (doc %d): %v", path, i+1, err)
			continue
		}

		if gvk.Kind == "" || gvk.Version == "" {
			continue
		}

		if !ctestglobals.IsIncludedKind(gvk.Kind) {
			continue
		}

		objects = append(objects, K8sObject{
			File:    path,
			Kind:    gvk.Kind,
			Name:    extractObjectName(obj),
			Object:  obj,
			RawYAML: doc,
		})
	}

	return objects, nil
}

// extractObjectName extracts the name from a Kubernetes object
func extractObjectName(obj runtime.Object) string {
	switch typed := obj.(type) {
	case interface{ GetName() string }:
		return typed.GetName()
	default:
		return "unknown"
	}
}
