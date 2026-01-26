package testrewrite

// ===============================
// Example 1 (Sysctl)
// ===============================

var OneShotUserExample = `
ORIGINAL FILE:

package node

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/kubernetes/test/e2e/environment"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
	imageutils "k8s.io/kubernetes/test/utils/image"
	admissionapi "k8s.io/pod-security-admission/api"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = SIGDescribe("Sysctls [LinuxOnly]", framework.WithNodeConformance(), func() {

	ginkgo.BeforeEach(func() {
		e2eskipper.SkipIfNodeOSDistroIs("windows")
	})

	f := framework.NewDefaultFramework("sysctl")
	f.NamespacePodSecurityLevel = admissionapi.LevelPrivileged
	var podClient *e2epod.PodClient

	testPod := func() *v1.Pod {
		podName := "sysctl-" + string(uuid.NewUUID())
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        podName,
				Annotations: map[string]string{},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "test-container",
						Image: imageutils.GetE2EImage(imageutils.BusyBox),
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		}
		return &pod
	}

	ginkgo.BeforeEach(func() {
		podClient = e2epod.NewPodClient(f)
	})

	framework.ConformanceIt("should support sysctls [MinimumKubeletVersion:1.21]", environment.NotInUserNS, func(ctx context.Context) {
		pod := testPod()
		pod.Spec.SecurityContext = &v1.PodSecurityContext{
			Sysctls: []v1.Sysctl{
				{
					Name:  "kernel.shm_rmid_forced",
					Value: "1",
				},
			},
		}
		pod.Spec.Containers[0].Command = []string{"/bin/sysctl", "kernel.shm_rmid_forced"}
	})

})
`

var OneShotAssistantExample = `
REWRITTEN FILE:

package node

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"

	ctest "k8s.io/kubernetes/test/ctest"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	ctestutils "k8s.io/kubernetes/test/ctest/utils"
)

var _ = SIGDescribe("Sysctls [LinuxOnly]", func() {

	framework.ConformanceIt("should support sysctls [MinimumKubeletVersion:1.21]", func(ctx context.Context) {

		fmt.Println(ctestglobals.StartSeparator)

		configs := getHardCodedConfigInfoSysctl()

		item, found := ctestutils.GetItemByExactTestInfo(configs, "default pod spec")
		if !found {
			framework.Failf("Get default hardcoded config failed.")
		}

		fmt.Println(ctestglobals.StartExtendModeSeparator)

		configObjs, _, err := ctest.GenerateEffectiveConfigReturnType[v1.PodSpec](item, ctest.ExtendOnly)
		if err != nil {
			framework.Failf("Failed to get matched fixtures: %v", err)
		}
		if configObjs != nil { 
			fmt.Println(ctestglobals.DebugPrefix(), "New Json Test Configs:", string(configJson)) 
			fmt.Println(ctestglobals.DebugPrefix(), "Num of Test Cases:", len(configObjs)) 
			fmt.Println("Start test config objs...")
			for i, configObj := range configObjs {
				fmt.Printf("Running %d th test cases.\n", i)
				fmt.Println(configObj)
				name := "sysctl-" + string(uuid.NewUUID())
				pod := &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec:       configObj,
				}

				pod.Spec.SecurityContext = &v1.PodSecurityContext{
					Sysctls: []v1.Sysctl{
						{Name: "kernel.shm_rmid_forced", Value: "1"},
					},
				}
			}
		}

		fmt.Println(ctestglobals.EndSeparator)
	})

})

func getHardCodedConfigInfoSysctl() ctestglobals.HardcodedConfig {
	return ctestglobals.HardcodedConfig{
		{
			FixtureFileName: "test_fixture.json",
			TestInfo:        []string{"default pod spec"},
			Field:           "spec",
			K8sObjects:      ctestglobals.PodSpecIncludeObjects,
			HardcodedConfig: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "test-container",
						Image: imageutils.GetE2EImage(imageutils.BusyBox),
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		},
	}
}
`

// ===============================
// Example 2 (ConfigMap)
// ===============================

var OneShotUserExample2 = `
ORIGINAL FILE:

package node

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epodoutput "k8s.io/kubernetes/test/e2e/framework/pod/output"
	imageutils "k8s.io/kubernetes/test/utils/image"
	admissionapi "k8s.io/pod-security-admission/api"

	"github.com/onsi/ginkgo/v2"
)

var _ = SIGDescribe("ConfigMap", func() {
	f := framework.NewDefaultFramework("configmap")
	f.NamespacePodSecurityLevel = admissionapi.LevelBaseline

	framework.ConformanceIt("should be consumable via the environment", f.WithNodeConformance(), func(ctx context.Context) {

		name := "configmap-test-" + string(uuid.NewUUID())
		configMap := newConfigMap(f, name)

		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pod-configmaps-" + string(uuid.NewUUID()),
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:    "env-test",
						Image:   imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{"sh", "-c", "env"},
						EnvFrom: []v1.EnvFromSource{
							{
								ConfigMapRef: &v1.ConfigMapEnvSource{
									LocalObjectReference: v1.LocalObjectReference{Name: name},
								},
							},
						},
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		}
	})
})
`

var OneShotAssistantExample2 = `
REWRITTEN FILE:

package node

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctest "k8s.io/kubernetes/test/ctest"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	ctestutils "k8s.io/kubernetes/test/ctest/utils"
)

var _ = SIGDescribe("ConfigMap", func() {

	framework.ConformanceIt("should be consumable via the environment", func(ctx context.Context) {

		fmt.Println(ctestglobals.StartSeparator)

		configs := getHardCodedConfigInfoConfigMap()

		item, found := ctestutils.GetItemByExactTestInfo(configs, "default configmap")
		if !found {
			framework.Failf("Get default hardcoded config failed.")
		}

		fmt.Println(ctestglobals.StartOverrideModeSeparator)

		configObjs, _, err := ctest.GenerateEffectiveConfigReturnType[map[string]string](item, ctest.OverrideOnly)
		if err != nil {
			framework.Failf("Failed to get matched fixtures: %v", err)
		}

		for _, data := range configObjs {
			_ = data
		}

		fmt.Println(ctestglobals.EndSeparator)
	})
})

func getHardCodedConfigInfoConfigMap() ctestglobals.HardcodedConfig {
	return ctestglobals.HardcodedConfig{
		{
			FixtureFileName: "test_fixture.json",
			TestInfo:        []string{"default configmap"},
			Field:           "data",
			K8sObjects:      []string{"configmaps"},
			HardcodedConfig: map[string]string{
				"data-1": "value-1",
				"data-2": "value-2",
				"data-3": "value-3",
			},
		},
	}
}
`
