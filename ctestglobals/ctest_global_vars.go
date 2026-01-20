package ctestglobals

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

// type Mode int

// const (
// 	ExtendOnly Mode = iota
// 	OverrideOnly
// 	Both
// )

const TestExternalFixtureFile = "test_fixtures.json"

type HardcodedConfig []struct {
	FixtureFileName string
	TestInfo        []string
	Field           string
	K8sObjects      []string
	HardcodedConfig interface{}
}

// HardcodedConfigItem describes one element of your HardcodedConfig slice.
type HardcodedConfigItem struct {
	FixtureFileName string
	TestInfo        []string
	Field           string
	K8sObjects      []string
	HardcodedConfig interface{}
}

var (
	StartSeparator             = "\n==================== CTEST START ===================="
	EndSeparator               = "\n==================== CTEST END ======================"
	StartExtendModeSeparator   = "\n==================== CTEST EXTEND ONLY START ===================="
	StartOverrideModeSeparator = "\n==================== CTEST OVERRIDE ONLY START ===================="
	StartUnionModeSeparator    = "\n==================== CTEST UNION MODE START ===================="
	KeyKind                    = "kind"
	KeyApiVersion              = "apiVersion"
	FixtureIncludeObjects      = []string{
		"Deployment",
		"StatefulSet",
		"DaemonSet",
		"ReplicaSet",
		"Pod",
		"Service",
		"ConfigMap",
		"Secret",
		"Namespace",
		"ServiceAccount",
		"PersistentVolume",
		"PersistentVolumeClaim",
		"ResourceQuota",
		"LimitRange",
		"Job",
		"CronJob",
		"Ingress",
		"NetworkPolicy",
		"Role",
		"RoleBinding",
		"ClusterRole",
		"ClusterRoleBinding",
		"StorageClass",
		"CustomResourceDefinition",
	}
	WeirdPaths            = []string{"github/workflows", ".github", ".travis.yml"}
	PodSpecIncludeObjects = []string{"deployments", "pods", "statefulSets", "daemonSets", "replicaSets"}
	DebugPrefix           = func() string {
		_, file, line, _ := runtime.Caller(1)
		return fmt.Sprintf("[DEBUG-CTEST %s %s:%d]:",
			time.Now().Format("2006-01-02 15:04:05"),
			file, line)
	}
	//configMap
	ConfigMapData = map[string]string{
		"data-1": "value-1",
		"data-2": "value-2",
		"data-3": "value-3",
	}
	// Global slices for keys and values
	ConfigMapKeys   []string
	ConfigMapValues []string

	//configMap
	SecretData = map[string][]byte{
		"data-1": []byte("value-1\n"),
		"data-2": []byte("value-2\n"),
		"data-3": []byte("value-3\n"),
	}

	SecretKeys   []string
	SecretValues []string

	//container probe
	LivenessProbe = &v1.Probe{
		InitialDelaySeconds: 15,
		TimeoutSeconds:      5,
		FailureThreshold:    1,
		PeriodSeconds:       30,
	}
	StartupProbe = &v1.Probe{
		InitialDelaySeconds: 15,
		TimeoutSeconds:      5,
		FailureThreshold:    3,
		PeriodSeconds:       5,
	}
	ReadinessProbe = &v1.Probe{
		InitialDelaySeconds: 15,
		TimeoutSeconds:      5,
		FailureThreshold:    1,
		PeriodSeconds:       2,
	}
	//downardapi
	Resources = v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("250m"),
			v1.ResourceMemory: resource.MustParse("32Mi"),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("1250m"),
			v1.ResourceMemory: resource.MustParse("64Mi"),
		},
	}

	ResourcesConverted = struct {
		LimitsCPU      int64
		LimitsMemory   int64
		RequestsCPU    int64
		RequestsMemory int64
	}{}

	Containers = []v1.Container{
		{
			Name:    "test-container-1",
			Image:   imageutils.GetE2EImage(imageutils.BusyBox),
			Command: []string{"/bin/sleep"},
			Args:    []string{"10000"},
			//pod level resources
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("150m"),
					v1.ResourceMemory: resource.MustParse("32Mi"),
				},
				Limits: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("250m"),
					v1.ResourceMemory: resource.MustParse("44Mi"),
				},
			},
			//privileged
			ImagePullPolicy: v1.PullIfNotPresent,
		},
	}
	//ephemeralContainer
	EphemeralContainer = v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{
			Name:    "debugger",
			Image:   imageutils.GetE2EImage(imageutils.BusyBox),
			Command: e2epod.GenerateScriptCmd("while true; do echo polo; sleep 2; done"),
			Stdin:   true,
			TTY:     true,
		},
	}

	hostnameOverride = "override.example.host"
	Spec             = v1.PodSpec{
		//host name override
		Hostname:         "custom-host",
		HostnameOverride: &hostnameOverride,
		Subdomain:        "t",
		Containers:       Containers,
		RestartPolicy:    v1.RestartPolicyNever,
		//pod level resources
		Resources: &v1.ResourceRequirements{
			Requests: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("250m"),
				v1.ResourceMemory: resource.MustParse("64Mi"),
			},
			Limits: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("1250m"),
				v1.ResourceMemory: resource.MustParse("128Mi"),
			},
		},
	}

	Spec1 = v1.PodSpec{
		//host name override
		Hostname:         "custom-host",
		HostnameOverride: &hostnameOverride,
		Subdomain:        "t",

		RestartPolicy: v1.RestartPolicyNever,
	}

	//pod resize
	// RR = struct {
	// 	LimitsCPU      string
	// 	LimitsMemory   string
	// 	RequestsCPU    string
	// 	RequestsMemory string
	// }{}

	PodTemplateSpec = v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{Name: "nginx", Image: imageutils.GetE2EImage(imageutils.Nginx)},
			},
		},
	}

	PodTemplate = v1.PodTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name: "podTemplateName",
			Labels: map[string]string{
				"podtemplate-static": "true", //label for list template, can not be changed
			},
		},
		Template: v1.PodTemplateSpec{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Name: "nginx", Image: imageutils.GetE2EImage(imageutils.Nginx)},
				},
			},
		},
	}
)

// Helper function
func MapKeysAndValues(m map[string]string) ([]string, []string) {
	keys := make([]string, 0, len(m))
	values := make([]string, 0, len(m))

	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

func extractKeysAndValues(secretData map[string][]byte) (keys []string, values []string) {
	for k, v := range secretData {
		keys = append(keys, k)
		values = append(values, string(v)) // convert []byte to string
	}
	return
}

func quantityToMilliCPU(q resource.Quantity) int64 {
	s := q.String() // e.g., "250m" or "1"
	if strings.HasSuffix(s, "m") {
		val, _ := strconv.ParseInt(strings.TrimSuffix(s, "m"), 10, 64)
		return val // already in millicores
	}
	val, _ := strconv.ParseFloat(s, 64)
	return int64(val * 1000) // convert cores to millicores
}

func quantityToBytes(q resource.Quantity) int64 {
	s := q.String() // e.g., "32Mi", "64Mi", "1Gi"
	// handle "Ki", "Mi", "Gi"
	multipliers := map[string]int64{
		"Ki": 1024,
		"Mi": 1024 * 1024,
		"Gi": 1024 * 1024 * 1024,
	}
	for suffix, mult := range multipliers {
		if strings.HasSuffix(s, suffix) {
			val, _ := strconv.ParseFloat(strings.TrimSuffix(s, suffix), 64)
			return int64(val * float64(mult))
		}
	}
	// fallback: plain bytes
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

// init function to initialize global key/value slices
func init() {
	ConfigMapKeys, ConfigMapValues = MapKeysAndValues(ConfigMapData)
	SecretKeys, SecretValues = extractKeysAndValues(SecretData)

	ResourcesConverted.RequestsCPU = quantityToMilliCPU(Resources.Requests[v1.ResourceCPU])
	ResourcesConverted.RequestsMemory = quantityToBytes(Resources.Requests[v1.ResourceMemory])
	ResourcesConverted.LimitsCPU = quantityToMilliCPU(Resources.Limits[v1.ResourceCPU])
	ResourcesConverted.LimitsMemory = quantityToBytes(Resources.Limits[v1.ResourceMemory])
}

func IsIncludedKind(kind string) bool {
	for _, k := range FixtureIncludeObjects {
		if k == kind {
			return true
		}
	}
	return false
}
