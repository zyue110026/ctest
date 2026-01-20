package ctest

import (
	// "encoding/json"
	"fmt"
	// "log"
	// "reflect"
	"testing"

	// fixtures "k8s.io/kubernetes/test/ctest/fixtures"
	v1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/api/resource"
	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
	utils "k8s.io/kubernetes/test/ctest/utils"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	// imageutils "k8s.io/kubernetes/test/utils/image"
	// k8sjson "k8s.io/apimachinery/pkg/util/json"
)

func getHardCodedConfigInfo() ctestglobals.HardcodedConfig {
	currentFile := utils.GetCurrentFileName()
	return ctestglobals.HardcodedConfig{{
		FixtureFileName: currentFile + "_fixture.json",
		TestInfo:        []string{"should be restarted with a exec \"cat /tmp/health\" liveness probe", "should be consumable via the environment", "should fail to create ConfigMap with empty key"},
		Field:           "containers",
		K8sObjects:      []string{""},
		HardcodedConfig: v1.PodSpec{
			//host name override
			Hostname: "custom-host",
			// HostnameOverride: &hostnameOverride,
			Subdomain:     "t",
			Containers:    ctestglobals.Containers,
			RestartPolicy: v1.RestartPolicyNever,
		},
	}}
}

func TestGenerateEffectiveConfig(t *testing.T) {
	// Get your configs
	configs := getHardCodedConfigInfo()
	// 1. Basic search
	item, found := utils.GetItemByExactTestInfoAndField(configs, "should be restarted with a exec \"cat /tmp/health\" liveness probe", "containers")
	if !found {
		t.Fatalf("Failed to find config item by TestInfo")
	}
	fmt.Println(item)
	// configObj, hardcodedType, configJson, nil := GenerateEffectiveConfig(item, Union)
	configObjs, configJson, nil := GenerateEffectiveConfigReturnType[v1.PodSpec](item, ExtendOnly)
	if nil != nil {
		t.Fatalf("Failed to get matched fixtures: %v", nil)
	}
	for _, configObj := range configObjs {
		fmt.Println(ctestglobals.DebugPrefix(), "Starting printing each configuration...")
		fmt.Println(configObj)
		// fmt.Println(string(configJson))
		// cmd := []string{"/bin/sh", "-c", "echo ok >/tmp/health; sleep 10; rm -rf /tmp/health; sleep 600"}
		// livenessProbe := configObj
		// pod := busyBoxPodSpec(nil, livenessProbe, cmd)
		// RunLivenessTest(ctx, f, pod, 1, defaultObservationTimeout)
	}

	fmt.Println(string(configJson))
	// fmt.Println(hardcodedType)

	// // Convert each json in configJson list to hardcodedType
	// for _, jsonData := range configJson {
	// 	var targetPtr reflect.Value
	// 	if hardcodedType.Kind() == reflect.Ptr {
	// 		targetPtr = reflect.New(hardcodedType.Elem())
	// 	} else {
	// 		targetPtr = reflect.New(hardcodedType)
	// 	}
	// 	if err := k8sjson.Unmarshal(jsonData, targetPtr.Interface()); err != nil {
	// 		fmt.Errorf("k8s json unmarshal into %s failed: %w", hardcodedType.String(), err)
	// 	}

	// 	var reconstructed reflect.Value
	// 	if targetPtr.Kind() == reflect.Ptr {
	// 		reconstructed = targetPtr.Elem()
	// 	} else {
	// 		reconstructed = targetPtr
	// 	}

	// 	unmarshaledObj := reconstructed.Interface()
	// 	fmt.Println(unmarshaledObj)
	// }
}
