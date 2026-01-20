package ctest

import (
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/kubernetes/test/ctest/fixtures"
)

// K8sObject represents one validated Kubernetes object decoded from YAML
type K8sObject struct {
	File    string
	Kind    string
	Name    string
	Object  interface{}
	RawYAML string
}

func ProcessObjects(objects []K8sObject) {
	for _, o := range objects {
		switch obj := o.Object.(type) {

		case *appsv1.Deployment:
			fixtures.AddDeployment(obj)

		case *appsv1.StatefulSet:
			fixtures.AddStatefulSet(obj)

		case *appsv1.DaemonSet:
			fixtures.AddDaemonSet(obj)

		case *appsv1.ReplicaSet:
			fixtures.AddReplicaSet(obj)

		case *batchv1.Job:
			fixtures.AddJob(obj)

		case *batchv1.CronJob:
			fixtures.AddCronJob(obj)

		case *corev1.Pod:
			fixtures.AddPod(obj)

		case *corev1.Service:
			fixtures.AddService(obj)

		case *corev1.ConfigMap:
			fixtures.AddConfigMap(obj)

		case *corev1.Secret:
			fixtures.AddSecret(obj)

		case *corev1.Namespace:
			fixtures.AddNamespace(obj)

		case *corev1.ServiceAccount:
			fixtures.AddServiceAccount(obj)

		case *corev1.PersistentVolume:
			fixtures.AddPersistentVolume(obj)

		case *corev1.PersistentVolumeClaim:
			fixtures.AddPersistentVolumeClaim(obj)

		case *corev1.ResourceQuota:
			fixtures.AddResourceQuota(obj)

		case *corev1.LimitRange:
			fixtures.AddLimitRange(obj)

		case *networkingv1.Ingress:
			fixtures.AddIngress(obj)

		case *networkingv1.NetworkPolicy:
			fixtures.AddNetworkPolicy(obj)

		case *rbacv1.Role:
			fixtures.AddRole(obj)

		case *rbacv1.RoleBinding:
			fixtures.AddRoleBinding(obj)

		case *rbacv1.ClusterRole:
			fixtures.AddClusterRole(obj)

		case *rbacv1.ClusterRoleBinding:
			fixtures.AddClusterRoleBinding(obj)

		case *storagev1.StorageClass:
			fixtures.AddStorageClass(obj)

		case *apiextv1.CustomResourceDefinition:
			fixtures.AddCustomResourceDefinition(obj)

		default:
			// Should not happen due to earlier filtering,
			// but we keep this to be defensive.
			meta, ok := obj.(metav1.Object)
			if ok {
				log.Printf(
					"skipping unsupported object kind=%s name=%s file=%s",
					o.Kind,
					meta.GetName(),
					o.File,
				)
			} else {
				log.Printf(
					"skipping unsupported object kind=%s file=%s",
					o.Kind,
					o.File,
				)
			}
		}
	}

	if err := fixtures.SaveFixtures(); err != nil {
		log.Fatalf("failed to save fixtures: %v", err)
	}

	// Print summary
	counts := fixtures.GetCounts()
	fmt.Printf("\n📊 Fixtures generated and saved:\n")
	total := 0
	for kind, count := range counts {
		if count > 0 {
			fmt.Printf("  %s: %d\n", kind, count)
			total += count
		}
	}

	if total == 0 {
		fmt.Printf("No fixtures were loaded from YAML files")
	} else {
		fmt.Printf("✅ Total: %d fixtures saved to file\n", total)
	}
}
