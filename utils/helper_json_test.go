package utils

import (
	// "encoding/json"
	"fmt"
	"log"
	"testing"

	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	fixtures "k8s.io/kubernetes/test/ctest/fixtures"
)

func TestGetFieldValuesFromFixtures(t *testing.T) {
	// 1) Load all non-null fixtures
	fixtures, err := fixtures.LoadFixturesAsJSON("../fixtures/" + ctestglobals.TestExternalFixtureFile)
	if err != nil {
		log.Fatalf("load all fixtures failed: %v", err)
	}
	// 1) find all top-level "containers" keys anywhere:
	vals, err := GetFieldValuesFromFixtures(fixtures, "allowPrivilegeEscalation")
	if err != nil {
		fmt.Println("err:", err)
	}
	PrintJSONRawMessages(vals)

	// 2) find exact path spec.template.spec.containers within deployments only:
	vals2, err := GetFieldValuesFromFixtures(fixtures, "spec.template.spec.containers.securityContext.allowPrivilegeEscalation", "deployments")
	if err != nil {
		fmt.Println("err:", err)
	}
	PrintJSONRawMessages(vals2)

}
