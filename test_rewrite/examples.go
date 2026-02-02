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

var OneShotUserExample3 = `
ORIGINAL GO Test Function:
func TestAdoption(t *testing.T) {
	testCases := []struct {
		name                    string
		existingOwnerReferences func(rs *apps.ReplicaSet) []metav1.OwnerReference
		expectedOwnerReferences func(rs *apps.ReplicaSet) []metav1.OwnerReference
	}{
		{
			"pod refers rs as an owner, not a controller",
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet"}}
			},
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true), BlockOwnerDeletion: ptr.To(true)}}
			},
		},
		{
			"pod doesn't have owner references",
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{}
			},
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true), BlockOwnerDeletion: ptr.To(true)}}
			},
		},
		{
			"pod refers rs as a controller",
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)}}
			},
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)}}
			},
		},
		{
			"pod refers other rs as the controller, refers the rs as an owner",
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{
					{UID: "1", Name: "anotherRS", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
					{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet"},
				}
			},
			func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{
					{UID: "1", Name: "anotherRS", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
					{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet"},
				}
			},
		},
	}
	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tCtx, closeFn, rm, informers, clientSet := rmSetup(t)
			defer closeFn()

			ns := framework.CreateNamespaceOrDie(clientSet, fmt.Sprintf("rs-adoption-%d", i), t)
			defer framework.DeleteNamespaceOrDie(clientSet, ns, t)

			rsClient := clientSet.AppsV1().ReplicaSets(ns.Name)
			podClient := clientSet.CoreV1().Pods(ns.Name)
			const rsName = "rs"
			rs, err := rsClient.Create(tCtx, newRS(rsName, ns.Name, 1), metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Failed to create replica set: %v", err)
			}
			podName := fmt.Sprintf("pod%d", i)
			pod := newMatchingPod(podName, ns.Name)
			pod.OwnerReferences = tc.existingOwnerReferences(rs)
			_, err = podClient.Create(tCtx, pod, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Failed to create Pod: %v", err)
			}

			stopControllers := runControllerAndInformers(t, rm, informers, 1)
			defer stopControllers()
			if err := wait.PollImmediate(interval, timeout, func() (bool, error) {
				updatedPod, err := podClient.Get(tCtx, pod.Name, metav1.GetOptions{})
				if err != nil {
					return false, err
				}

				e, a := tc.expectedOwnerReferences(rs), updatedPod.OwnerReferences
				if reflect.DeepEqual(e, a) {
					return true, nil
				}

				t.Logf("ownerReferences don't match, expect %v, got %v", e, a)
				return false, nil
			}); err != nil {
				t.Fatalf("test %q failed: %v", tc.name, err)
			}
		})
	}
}
`
var OneShotAssistantExample3 = `
REWRITTEN GO Test FUNCTION:
// Rewritten TestAdoption with edge/invalid test cases
func TestCtestAdoptionEdgeCases(t *testing.T) {
	
	edgeTestCases := []struct {
		name                    string
		existingOwnerReferences func(rs *apps.ReplicaSet) []metav1.OwnerReference
		expectedOwnerReferences func(rs *apps.ReplicaSet) []metav1.OwnerReference
	}{
		{
			name: "pod has multiple controller owners",
			existingOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{
					{UID: "1", Name: "rs1", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
					{UID: "2", Name: "rs2", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
				}
			},
			expectedOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				// No adoption should happen since pod already has a controller
				return []metav1.OwnerReference{
					{UID: "1", Name: "rs1", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
					{UID: "2", Name: "rs2", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
				}
			},
		},
		{
			name: "pod has owner reference with invalid UID",
			existingOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{
					{UID: "", Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet"},
				}
			},
			expectedOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				// Pod should not adopt due to invalid UID
				return []metav1.OwnerReference{
					{UID: "", Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet"},
				}
			},
		},
		{
			name: "pod has owner reference with unknown kind",
			existingOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{
					{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "UnknownKind"},
				}
			},
			expectedOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				// Should not adopt since kind is not ReplicaSet
				return []metav1.OwnerReference{
					{UID: rs.UID, Name: rs.Name, APIVersion: "apps/v1", Kind: "UnknownKind"},
				}
			},
		},
		{
			name: "pod has owner reference with invalid APIVersion",
			existingOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{
					{UID: rs.UID, Name: rs.Name, APIVersion: "invalid/v1", Kind: "ReplicaSet"},
				}
			},
			expectedOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				// Adoption should fail due to wrong APIVersion
				return []metav1.OwnerReference{
					{UID: rs.UID, Name: rs.Name, APIVersion: "invalid/v1", Kind: "ReplicaSet"},
				}
			},
		},
		{
			name: "pod has multiple owners including invalid and correct controller",
			existingOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				return []metav1.OwnerReference{
					{UID: "random", Name: "otherRS", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
					{UID: "", Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet"},
				}
			},
			expectedOwnerReferences: func(rs *apps.ReplicaSet) []metav1.OwnerReference {
				// Should not adopt due to existing controller present
				return []metav1.OwnerReference{
					{UID: "random", Name: "otherRS", APIVersion: "apps/v1", Kind: "ReplicaSet", Controller: ptr.To(true)},
					{UID: "", Name: rs.Name, APIVersion: "apps/v1", Kind: "ReplicaSet"},
				}
			},
		},
	}
	fmt.Println(ctestglobals.DebugPrefix(), "Add edge test cases:", edgeTestCases)
	fmt.Println(ctestglobals.DebugPrefix(), "Number of test cases:", len(edgeTestCases))
	for i, tc := range edgeTestCases {
		fmt.Printf("Running %d th test cases.\n", i)
		fmt.Println(tc)
		t.Run(tc.name, func(t *testing.T) {
			tCtx, closeFn, rm, informers, clientSet := rmSetup(t)
			defer closeFn()

			ns := framework.CreateNamespaceOrDie(clientSet, fmt.Sprintf("rs-adoption-edge-%d", i), t)
			defer framework.DeleteNamespaceOrDie(clientSet, ns, t)

			rsClient := clientSet.AppsV1().ReplicaSets(ns.Name)
			podClient := clientSet.CoreV1().Pods(ns.Name)
			rsName := fmt.Sprintf("rs-%s", string(uuid.NewUUID()))
			rs, err := rsClient.Create(tCtx, newRS(rsName, ns.Name, 1), metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Failed to create replica set: %v", err)
			}

			podName := fmt.Sprintf("pod-edge-%d", i)
			pod := newMatchingPod(podName, ns.Name)
			pod.OwnerReferences = tc.existingOwnerReferences(rs)
			_, err = podClient.Create(tCtx, pod, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Failed to create Pod: %v", err)
			}

			stopControllers := runControllerAndInformers(t, rm, informers, 1)
			defer stopControllers()

			if err := wait.PollImmediate(interval, timeout, func() (bool, error) {
				updatedPod, err := podClient.Get(tCtx, pod.Name, metav1.GetOptions{})
				if err != nil {
					return false, err
				}

				e, a := tc.expectedOwnerReferences(rs), updatedPod.OwnerReferences
				if reflect.DeepEqual(e, a) {
					return true, nil
				}

				t.Logf("ownerReferences don't match, expect %v, got %v", e, a)
				return false, nil
			}); err != nil {
				t.Fatalf("edge test %q failed: %v", tc.name, err)
			}
		})
	}
}
`
