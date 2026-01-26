package testrewrite

import (
	"fmt"
	"path/filepath"
	// "strings"
)

var FixtureIncludeObjects = []string{
	"deployments",
	"statefulsets",
	"daemonsets",
	"replicasets",
	"pods",
	"services",
	"configmaps",
	"secrets",
	"namespaces",
	"serviceaccounts",
	"persistentvolumes",
	"persistentvolumeclaims",
	"resourcequotas",
	"limitranges",
	"jobs",
	"cronjob",
	"ingressws",
	"networkpolicys",
	"roles",
	"rolebindings",
	"clusterroles",
	"clusterrolebindings",
	"storageclasses",
	"customresourcedefinitions",
}

// BuildPrompt generates the DeepSeek-Coder prompt for rewriting a Kubernetes test file
func BuildPrompt(path, content string) string {
	fileName := filepath.Base(path)

	// Prompt template
	prompt := fmt.Sprintf(`
You are a Go developer rewriting Kubernetes tests for dynamic configuration.

The goal of rewriting is to remove hard-coded Kubernetes configuration values from tests and replace them with dynamically generated configuration, while preserving the original intent and semantics of the test.

Kubernetes tests often validate behavior under specific configuration assumptions (for example, restartPolicy: Never). Not all configuration fields are exercised in a single test. We want to evaluate whether the test still succeeds when configuration is dynamically:
1) Extended with additional fields,
2) Overridden with different values,
3) Or both extended and overridden.

However, rewritten tests MUST NOT break the original test logic or change what the test is intended to verify.

Before rewriting, you must carefully analyze the original test code and understand:
- What behavior the test is validating
- Which configuration fields are essential to the test's correctness
- Which fields are merely incidental and safe to vary

Instructions:

1. **Package and imports**:
   - Keep the original package.
   - Keep needed imports.
   - Add imports if missing:
     import (
         "fmt"
         ctest "k8s.io/kubernetes/test/ctest"
         ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
         ctestutils "k8s.io/kubernetes/test/ctest/utils"
     )

2. **Hardcoded Config Function**:
   - For each hardcoded config in a test, generate:
   
     func getHardCodedConfigInfo<FileName>() ctestglobals.HardcodedConfig
     
     - Returns:
       
       type HardcodedConfig []struct {
           FixtureFileName string
           TestInfo        []string
           Field           string
           K8sObjects      []string
           HardcodedConfig interface{}
       }
       
     - Populate:
       - FixtureFileName: "test_fixture.json"
       - TestInfo: unique test description string
       - Field: field name mapping to hardcoded values
       - K8sObjects: choose all relevant objects from "FixtureIncludeObjects"
       - HardcodedConfig: exact hardcoded values from original test (only the part necessary for the test)
	 - Example structure for container_probe.go:
     func getHardCodedConfigInfoContainerProbe() ctestglobals.HardcodedConfig {
         return ctestglobals.HardcodedConfig{
			{
				FixtureFileName: "test_fixture.json",
				TestInfo: []string{
					"should be restarted with a local redirect http liveness probe"},
				Field:      "livenessProbe",
				K8sObjects: []string{"deployments", "pods", "statefulSets", "daemonSets", "replicaSets"},
				HardcodedConfig: &v1.Probe{
					ProbeHandler:        httpGetHandler("/redirect?loc="+url.QueryEscape("/healthz"), 8080),
					InitialDelaySeconds: 15,
					FailureThreshold:    1,
				}, 
			}
         }
     }
	

3. **Hardcoded Config Selection**:
   - Only store the minimum part of the object needed for the test
     (e.g., PodSpec instead of entire Pod if testing security context).
   - K8sObjects must be selected from:
     FixtureIncludeObjects = []string{
       "deployments","statefulsets","daemonsets","replicasets","pods","services","configmaps",
       "secrets","namespaces","serviceaccounts","persistentvolumes","persistentvolumeclaims",
       "resourcequotas","limitranges","jobs","cronjob","ingressws","networkpolicys",
       "roles","rolebindings","clusterroles","clusterrolebindings","storageclasses","customresourcedefinitions"
     }
   - Select all relevant objects that could contain the hardcoded field. For example, if the hardcoded field is in Spec, include "pods", "deployments", "statefulsets", "daemonsets", "replicasets"...

4. **Rewriting Tests**:
   - Preserve all dynamic fields and metadata.
   - Replace only the hardcoded config with generated configObjs from:
     configObjs, configJson, err := ctest.GenerateEffectiveConfigReturnType[<type>](item, <mode>)
   - Inject dynamic or predefined configuration values, for example:
     name := "<prefix>-" + string(uuid.NewUUID())
     configObj.Containers[0].Name = name
     pod := &v1.Pod{
         ObjectMeta: metav1.ObjectMeta{Name: name},
         Spec:       configObj, // replace only the needed part
     }
   - Call the original test execution function with the new object.
   - Modify the original test execution function to accept the new configObj and make sure test purpose keeps same, if needed.

5. **Merge Mode Logic**:
   - Decide mode based on test safety:
     - Only extend: ctest.ExtendOnly, use ctestglobals.StartExtendModeSeparator
     - Override only: ctest.OverrideOnly, use ctestglobals.StartOverrideModeSeparator
     - Union: ctest.Union, use ctestglobals.StartUnionModeSeparator
   - Print the separator before starting the rewritten test.

6. **Logging and Debug**:
   - Use fmt.Println(ctestglobals.DebugPrefix(), "message") for logging.
   - Always log:
     - Start of test
     - Matched config
     - JSON of new test config
     - Number of test cases
   - Handle errors using framework.Failf if config not found.

7. **Preserve Original Functions**:
   - Keep all other helper functions in the original file.
   - Only add new test functions and getHardCodedConfigInfo functions.

8. **Multiple Tests per File**:
   - Repeat process for each test case.
   - Each test must have unique TestInfo.

9. **Output**:
   - If this file has tests that need rewriting, return the code exactly.
   - Ensure the file compiles and runs.
   - Do not include any explanations or comments outside the code.
   - If this file has no tests that need rewriting, return the string "NONE" exactly.

Below is original go code content, generate the rewritten Go test code based on the above instructions and below original code. 
---
Input file: %s

File content:
%s

---
`, fileName, content)

	return prompt
}
