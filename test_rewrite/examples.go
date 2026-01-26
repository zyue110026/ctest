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
		// sysctl is not supported on Windows.
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

	/*
	  Release: v1.21
	  Testname: Sysctl, test sysctls
	  Description: Pod is created with kernel.shm_rmid_forced sysctl. Kernel.shm_rmid_forced must be set to 1
	  [LinuxOnly]: This test is marked as LinuxOnly since Windows does not support sysctls
	*/
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

		ginkgo.By("Creating a pod with the kernel.shm_rmid_forced sysctl")
		pod = podClient.Create(ctx, pod)

		ginkgo.By("Watching for error events or started pod")
		// watch for events instead of termination of pod because the kubelet deletes
		// failed pods without running containers. This would create a race as the pod
		// might have already been deleted here.
		ev, err := e2epod.NewPodClient(f).WaitForErrorEventOrSuccess(ctx, pod)
		framework.ExpectNoError(err)
		gomega.Expect(ev).To(gomega.BeNil())

		ginkgo.By("Waiting for pod completion")
		err = e2epod.WaitForPodNoLongerRunningInNamespace(ctx, f.ClientSet, pod.Name, f.Namespace.Name)
		framework.ExpectNoError(err)
		pod, err = podClient.Get(ctx, pod.Name, metav1.GetOptions{})
		framework.ExpectNoError(err)

		ginkgo.By("Checking that the pod succeeded")
		gomega.Expect(pod.Status.Phase).To(gomega.Equal(v1.PodSucceeded))

		ginkgo.By("Getting logs from the pod")
		log, err := e2epod.GetPodLogs(ctx, f.ClientSet, f.Namespace.Name, pod.Name, pod.Spec.Containers[0].Name)
		framework.ExpectNoError(err)

		ginkgo.By("Checking that the sysctl is actually updated")
		gomega.Expect(log).To(gomega.ContainSubstring("kernel.shm_rmid_forced = 1"))
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

	"k8s.io/kubernetes/test/e2e/environment"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
	imageutils "k8s.io/kubernetes/test/utils/image"
	admissionapi "k8s.io/pod-security-admission/api"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	ctest "k8s.io/kubernetes/test/ctest"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	ctestutils "k8s.io/kubernetes/test/ctest/utils"
)

var _ = SIGDescribe("Sysctls [LinuxOnly]", framework.WithNodeConformance(), func() {

	ginkgo.BeforeEach(func() {
		// sysctl is not supported on Windows.
		e2eskipper.SkipIfNodeOSDistroIs("windows")
	})

	f := framework.NewDefaultFramework("sysctl")
	f.NamespacePodSecurityLevel = admissionapi.LevelPrivileged
	var podClient *e2epod.PodClient

	ginkgo.BeforeEach(func() {
		podClient = e2epod.NewPodClient(f)
	})

	/*
	  Release: v1.21
	  Testname: Sysctl, test sysctls
	  Description: Pod is created with kernel.shm_rmid_forced sysctl. Kernel.shm_rmid_forced must be set to 1
	  [LinuxOnly]: This test is marked as LinuxOnly since Windows does not support sysctls
	*/

	framework.ConformanceIt("should support sysctls [MinimumKubeletVersion:1.21]", environment.NotInUserNS, func(ctx context.Context) {
		fmt.Println(ctestglobals.StartSeparator)
		configs := getHardCodedConfigInfoSysctl()

		// 1. Basic search
		item, found := ctestutils.GetItemByExactTestInfo(configs, "default pod spec")
		if !found {
			fmt.Println(ctestglobals.DebugPrefix(), "Failed to find config item by TestInfo")
			framework.Failf("Get default hardcoded config failed.")
		}
		fmt.Println(ctestglobals.DebugPrefix(), "get default configs:", item)
		fmt.Println(ctestglobals.StartExtendModeSeparator)
		configObjs, configJson, err := ctest.GenerateEffectiveConfigReturnType[v1.PodSpec](item, ctest.ExtendOnly)
		if err != nil {
			fmt.Println(ctestglobals.DebugPrefix(), "Failed to get matched fixtures: %v", err)
			framework.Failf("Failed to get matched fixtures: %v", err)
		}
		if configObjs != nil {
			fmt.Println(ctestglobals.DebugPrefix(), "New Json Test Configs:", string(configJson))
			fmt.Println(ctestglobals.DebugPrefix(), "Num of Test Cases:", len(configObjs))
			fmt.Println("Start test config objs...")
			for i, configObj := range configObjs {
				fmt.Printf("Running %d th test cases.\n", i)
				fmt.Println(configObj)
				testPod := func() *v1.Pod {
					podName := "sysctl-" + string(uuid.NewUUID())
					pod := v1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:        podName,
							Annotations: map[string]string{},
						},
						Spec: configObj,
					}

					return &pod
				}
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

				ginkgo.By("Creating a pod with the kernel.shm_rmid_forced sysctl")
				pod = podClient.Create(ctx, pod)

				ginkgo.By("Watching for error events or started pod")
				// watch for events instead of termination of pod because the kubelet deletes
				// failed pods without running containers. This would create a race as the pod
				// might have already been deleted here.
				ev, err := e2epod.NewPodClient(f).WaitForErrorEventOrSuccess(ctx, pod)
				framework.ExpectNoError(err)
				gomega.Expect(ev).To(gomega.BeNil())

				ginkgo.By("Waiting for pod completion")
				err = e2epod.WaitForPodNoLongerRunningInNamespace(ctx, f.ClientSet, pod.Name, f.Namespace.Name)
				framework.ExpectNoError(err) //failed container test-container failed reason: container not ready.
				pod, err = podClient.Get(ctx, pod.Name, metav1.GetOptions{})
				framework.ExpectNoError(err)

				ginkgo.By("Checking that the pod succeeded")
				gomega.Expect(pod.Status.Phase).To(gomega.Equal(v1.PodSucceeded))

				ginkgo.By("Getting logs from the pod")
				log, err := e2epod.GetPodLogs(ctx, f.ClientSet, f.Namespace.Name, pod.Name, pod.Spec.Containers[0].Name)
				framework.ExpectNoError(err)

				ginkgo.By("Checking that the sysctl is actually updated")
				gomega.Expect(log).To(gomega.ContainSubstring("kernel.shm_rmid_forced = 1"))

			}
		} else {
			fmt.Println(ctestglobals.DebugPrefix(), "Skipping test execution. No new config objs found. ")
		}
		fmt.Println(ctestglobals.EndSeparator)

	})

})

func getHardCodedConfigInfoSysctl() ctestglobals.HardcodedConfig {
	return ctestglobals.HardcodedConfig{{
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
	}}
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

	/*
		Release: v1.9
		Testname: ConfigMap, from environment variables
		Description: Create a Pod with a environment source from ConfigMap. All ConfigMap values MUST be available as environment variables in the container.
	*/
	framework.ConformanceIt("should be consumable via the environment", f.WithNodeConformance(), func(ctx context.Context) {
		name := "configmap-test-" + string(uuid.NewUUID())
		configMap := newConfigMap(f, name)
		ginkgo.By(fmt.Sprintf("Creating configMap %v/%v", f.Namespace.Name, configMap.Name))
		var err error
		if configMap, err = f.ClientSet.CoreV1().ConfigMaps(f.Namespace.Name).Create(ctx, configMap, metav1.CreateOptions{}); err != nil {
			framework.Failf("unable to create test configMap %s: %v", configMap.Name, err)
		}

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
								ConfigMapRef: &v1.ConfigMapEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: name}},
							},
							{
								Prefix:       "p-",
								ConfigMapRef: &v1.ConfigMapEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: name}},
							},
						},
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		}

		e2epodoutput.TestContainerOutput(ctx, f, "consume configMaps", pod, 0, []string{
			"data-1=value-1", "data-2=value-2", "data-3=value-3",
			"p-data-1=value-1", "p-data-2=value-2", "p-data-3=value-3",
		})
	})

})

func newConfigMap(f *framework.Framework, name string) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: f.Namespace.Name,
			Name:      name,
		},
		Data: map[string]string{
			"data-1": "value-1",
			"data-2": "value-2",
			"data-3": "value-3",
		},
	}
}

func newConfigMapWithEmptyKey(ctx context.Context, f *framework.Framework) (*v1.ConfigMap, error) {
	name := "configmap-test-emptyKey-" + string(uuid.NewUUID())
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: f.Namespace.Name,
			Name:      name,
		},
		Data: map[string]string{
			"": "value-1",
		},
	}

	ginkgo.By(fmt.Sprintf("Creating configMap that has name %s", configMap.Name))
	return f.ClientSet.CoreV1().ConfigMaps(f.Namespace.Name).Create(ctx, configMap, metav1.CreateOptions{})
}

`

var OneShotAssistantExample2 = `
REWRITTEN FILE:

package node

import (
	"context"
	// "encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"

	"k8s.io/kubernetes/test/e2e/framework"
	e2epodoutput "k8s.io/kubernetes/test/e2e/framework/pod/output"
	imageutils "k8s.io/kubernetes/test/utils/image"
	admissionapi "k8s.io/pod-security-admission/api"

	"github.com/onsi/ginkgo/v2"
	//"github.com/onsi/gomega"

	ctest "k8s.io/kubernetes/test/ctest"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	ctestutils "k8s.io/kubernetes/test/ctest/utils"
)

var _ = SIGDescribe("ConfigMap", func() {
	f := framework.NewDefaultFramework("configmap")
	f.NamespacePodSecurityLevel = admissionapi.LevelBaseline

	/*
		Release: v1.9
		Testname: ConfigMap, from environment variables
		Description: Create a Pod with a environment source from ConfigMap. All ConfigMap values MUST be available as environment variables in the container.
	*/
	framework.ConformanceIt("should be consumable via the environment", f.WithNodeConformance(), func(ctx context.Context) {
		fmt.Println(ctestglobals.StartSeparator)
		configMapDatas, e := getConfigMapFromFixtureOverrideMode("default configmap")
		if e != nil {
			framework.Failf("Get configMap from fixture failed: %v", e)
		}
		if configMapDatas != nil {
			fmt.Println(ctestglobals.DebugPrefix(), "New Json Test Configs:", configMapDatas)
			fmt.Println(ctestglobals.DebugPrefix(), "Num of Test Cases:", len(configMapDatas))
			fmt.Println("Start test config objs...")
			for i, configMapData := range configMapDatas {
				configMapKeys, configMapValues := ctestglobals.MapKeysAndValues(configMapData)
				fmt.Printf("Running %d th test cases.\n", i)
				fmt.Println("ConfigMap Data:", configMapData)
				fmt.Println("ConfigMap Data Keys:", configMapKeys)
				fmt.Println("ConfigMap Data Values:", configMapValues)
				name := "configmap-test-" + string(uuid.NewUUID())
				configMap := newConfigMap(f, name)
				configMap.Data = configMapData
				ginkgo.By(fmt.Sprintf("Creating configMap %v/%v", f.Namespace.Name, configMap.Name))
				var err error
				if configMap, err = f.ClientSet.CoreV1().ConfigMaps(f.Namespace.Name).Create(ctx, configMap, metav1.CreateOptions{}); err != nil {
					framework.Failf("unable to create test configMap %s: %v", configMap.Name, err)
				}

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
										ConfigMapRef: &v1.ConfigMapEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: name}},
									},
									{
										Prefix:       "p-",
										ConfigMapRef: &v1.ConfigMapEnvSource{LocalObjectReference: v1.LocalObjectReference{Name: name}},
									},
								},
							},
						},
						RestartPolicy: v1.RestartPolicyNever,
					},
				}

				e2epodoutput.TestContainerOutput(ctx, f, "consume configMaps", pod, 0, func() []string {
					out := make([]string, 0, len(configMapKeys)*2)
					// no prefix
					for i := range configMapKeys {
						out = append(out, fmt.Sprintf("%s=%s", configMapKeys[i], configMapValues[i]))
					}
					// prefix "p-"
					for i := range configMapKeys {
						out = append(out, fmt.Sprintf("p-%s=%s", configMapKeys[i], configMapValues[i]))
					}
					return out
				}())

			}
		} else {
			fmt.Println(ctestglobals.DebugPrefix(), "Skipping test execution. No new config objs found. ")
		}
		fmt.Println(ctestglobals.EndSeparator)

	})

})

func getHardCodedConfigInfoConfigMap() ctestglobals.HardcodedConfig {
	return ctestglobals.HardcodedConfig{{
		FixtureFileName: "test_fixture.json",
		TestInfo:        []string{"default configmap"},
		Field:           "data",
		K8sObjects:      []string{"configmaps"},
		HardcodedConfig: map[string]string{
			"data-1": "value-1",
			"data-2": "value-2",
			"data-3": "value-3",
		},
	}}
}

func getConfigMapFromFixtureOverrideMode(testinfo string) ([]map[string]string, error) {
	hardcodedConfig := getHardCodedConfigInfoConfigMap()
	// 1. Basic search
	item, found := ctestutils.GetItemByExactTestInfo(hardcodedConfig, testinfo)
	if !found {
		fmt.Println(ctestglobals.DebugPrefix(), "Failed to find config item by TestInfo")
		framework.Failf("Get default hardcoded config failed.")
	}
	fmt.Println(ctestglobals.DebugPrefix(), "get default configs:", item)
	// fmt.Println(item)
	fmt.Println(ctestglobals.StartOverrideModeSeparator)
	configObjs, configJson, err := ctest.GenerateEffectiveConfigReturnType[map[string]string](item, ctest.OverrideOnly)
	if err != nil {
		fmt.Println(ctestglobals.DebugPrefix(), "Failed to get matched fixtures: %v", err)
		framework.Failf("Failed to get matched fixtures: %v", err)
	}
	if configObjs != nil {
		fmt.Println(ctestglobals.DebugPrefix(), "New Json Test Configs:", string(configJson))
		fmt.Println(ctestglobals.DebugPrefix(), "Num of Test Cases:", len(configObjs))

		return configObjs, nil
	} else {

		return nil, nil
	}

}

`
