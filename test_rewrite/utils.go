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

Kubernetes tests typically define testcases = [] to cover common scenarios, but they may not include edge cases or invalid values. In these situations, we should add test cases with edge conditions and invalid values to ensure the test behavior is still correct and to identify potential gaps in validation logic.

However, rewritten tests MUST NOT break the original test logic or change what the test is intended to verify.

Before rewriting, you must carefully analyze the original test code and understand:
- What behavior the test is validating
- Which configuration fields are essential to the test's correctness
- Which fields are merely incidental and safe to vary

Instructions:

1. **Package and imports**:
   - Keep the original package.
   - Keep needed imports.
   - Add imports as needed if missing:
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
       - Field: field name mapping to hardcoded values. *MUST macth exactly K8s object field name, e.g., "restartPolicy", "securityContext", "livenessProbe".*
       - K8sObjects: choose all relevant objects from "FixtureIncludeObjects"
       - HardcodedConfig: exact hardcoded values from original test (only the part necessary for the test). Do NOT include variables.
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
   - Select ALL relevant objects that could contain the hardcoded field. *For example, if the hardcoded field is in spec, you should select all relevant objects include "pods", "deployments", "statefulsets", "daemonsets", "replicasets", instead of just "pods".*

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

6. **Handling Test Cases**:
   - If the original test has testcases = []:
      - Add edge cases and invalid values (empty strings, nil pointers, zero, negative, extremely large values)
      - Preserve original test semantics, and add comments if needed to explain the purpose of edge cases.
   - If both hardcoded values and testcases exist, do both:
      - Generate dynamic configurations for hardcoded values using merge modes
      - Expand testcases array to include edge and invalid values
      - Run the test by combining all dynamic configs with all testcases
   - *If testcases exist and hardcpded cpmfiguration not related to k8s configuration field, only do add edge cases. If no testcases and hardcoded configuration exists, do not rewrite the test, and do not include it in the new file.*

7. **Logging and Debug**:
   - Use fmt.Println(ctestglobals.DebugPrefix(), "message") for logging.
   - Always log:
     - Start of test
     - Matched config, for example: fmt.Println(ctestglobals.DebugPrefix(), "get default configs:", item)
     - JSON of new test configs, for example: fmt.Println(ctestglobals.DebugPrefix(), "New Json Test Configs:", string(configJson))
     - Number of test cases, for example: fmt.Println(ctestglobals.DebugPrefix(), "Number of test cases:", len(configObjs))
     - For each test case, log the index and the config used. For example: fmt.Sprintf("Running # th test cases.\n", i)
				fmt.Println(configObj)
     - Skipped tests due to missing config, for example: fmt.Println(ctestglobals.DebugPrefix(), "Skipping test execution. No new configurations generated. "). Note, use if-else to check if configObjs is nil or empty. If configObjs == nil, skip the test execution, and simply log the skip message and continue run tests, do not use framework.Failf break test execution.
   - Handle errors using framework.Failf for ginkgo test, using t.Fatalf for go test function.

8. **Unchanged functions should never appear in the new file**: 
    - Only add new test functions, new helper functions, and getHardCodedConfigInfo<FileName> functions in the new file.
    - Each new helper function MUST has a unique name (for example, append <FileName>).
    - Only include functions that are new or modified for dynamic configuration and edge-case testing. Do not copy unchanged helpers or tests.

9. **Multiple Tests per File**:
   - For each test function provided, such as: framework.ConformanceIt, ginkgo.It, f.It, framework.It, ginkgo.Describe, framework.Describe, f.Describe, func TestXYZ etc. do:
    - Repeat process for each test function.
    - If you think a test function cannot be rewritten, do not include it.
    - Only return successfully rewritten tests.
    - Each test must have unique TestInfo in func getHardCodedConfigInfo<FileName>() ctestglobals.HardcodedConfig
    - Do NOT omit any part of the code for brevity within a rewritten test.
    - *If you are rewriting a go test function (func TestXYZ(t *testing.T)):
      - Make sure to rename it to append Ctest bewteen Test and original name, e.g., func TestCtestXYZ(t *testing.T), to avoid name conflicts with the original test function.*
      - Add some test cases for edge cases and log debug info, for example, empty values, max values, min values, etc., invalide value, if applicable.*
   - For all successfully rewritten tests in a file:
    - Collect all hardcoded configurations into one function: func getHardCodedConfigInfo<FileName>() ctestglobals.HardcodedConfig
    - Each entry in the returned HardcodedConfig slice corresponds to one test and contains its unique TestInfo.

10. **Output**:
   - If this file has tests that need rewriting, RETURN THE FINAL CODE EXACTLY.
   - Ensure the file compiles and runs, and remove all decleared but unused var and imports.
   - Do NOT include any explanations or comments outside the code.
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

// 7. **Helper Function Reuse and Modification**:
// The generated new file MUST NOT contain any helper function whose implementation is identical to one in the original file.
// If a helper function is unchanged, the rewritten test MUST call the original function by name and MUST NOT redefine it.
// Only include a helper function in the new file when its logic is changed to support dynamic configuration.
// Any modified helper function MUST:
//   - Have a unique name (for example, append Rewritten)
//   - Be used by the rewritten test
// If a helper function is unchanged, it is an error to include it in the new file.
