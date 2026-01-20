//go:build ctest
// +build ctest

package ctest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	//"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	fixtures "k8s.io/kubernetes/test/ctest/fixtures"
)

// K8sObject represents a parsed Kubernetes object with its original file
type K8sObject struct {
	File    string
	Kind    string
	Name    string
	Object  runtime.Object
	RawYAML string
}

// Config holds the application configuration
type Config struct {
	YAMLDir  string
	FileName string // if empty, process all files
	Verbose  bool
}

func TestGenerateFixtures(t *testing.T) {
	// Clear existing fixtures
	// fixtures.InitializeFixtures()
	fixtures.ClearFixtures()
	// Example usage
	config := Config{
		YAMLDir:  "./yamls",
		FileName: "", // empty = all files, or specify "deployment.yaml"
		Verbose:  true,
	}

	objects, err := ProcessYAMLFiles(config)
	if err != nil {
		log.Fatal(err)
	}

	// Print summary
	fmt.Printf("\n📊 Processed %d Kubernetes objects:\n", len(objects))
	for _, obj := range objects {
		fmt.Printf("  - %s: %s (%s)\n", obj.Kind, obj.Name, filepath.Base(obj.File))
	}

	// Demonstrate working with specific types
	ProcessObjects(objects)
}

// ProcessYAMLFiles reads and parses YAML files based on configuration
func ProcessYAMLFiles(config Config) ([]K8sObject, error) {
	var allObjects []K8sObject

	files, err := getYAMLFiles(config.YAMLDir, config.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get YAML files: %w", err)
	}

	if config.Verbose {
		fmt.Printf("🔍 Found %d YAML files to process:\n", len(files))
		for _, file := range files {
			fmt.Printf("   - %s\n", file)
		}
	}

	for _, file := range files {
		if config.Verbose {
			fmt.Printf("\n📄 Processing file: %s\n", file)
		}

		objects, err := processFile(file)
		if err != nil {
			log.Printf("⚠️ Warning: Failed to process file %s: %v", file, err)
			continue
		}

		allObjects = append(allObjects, objects...)
	}

	return allObjects, nil
}

// getYAMLFiles returns list of YAML files to process
func getYAMLFiles(dirPath, fileName string) ([]string, error) {
	if fileName != "" {
		// Single file mode
		fullPath := filepath.Join(dirPath, fileName)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", fullPath)
		}
		return []string{fullPath}, nil
	}

	// All files mode
	var yamlFiles []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})

	return yamlFiles, err
}

// processFile reads a single file and parses all YAML documents
func processFile(filePath string) ([]K8sObject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return parseMultiDocYAML(string(data), filePath)
}

// parseMultiDocYAML parses a multi-document YAML string
func parseMultiDocYAML(yamlContent, filePath string) ([]K8sObject, error) {
	var objects []K8sObject

	// Split the YAML by the document separator
	documents := strings.Split(yamlContent, "---")

	for docIndex, doc := range documents {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		if err := registerCustomSchemes(); err != nil {
			return nil, fmt.Errorf("failed to register custom schemes: %w", err)
		}

		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, gvk, err := decode([]byte(doc), nil, nil)
		if err != nil {
			log.Printf("⚠️ Failed to decode document %d in %s: %v", docIndex+1, filePath, err)
			continue
		}

		k8sObj := K8sObject{
			File:    filePath,
			Kind:    gvk.Kind,
			Object:  obj,
			RawYAML: doc,
		}

		// Extract name based on type
		k8sObj.Name = extractObjectName(obj)

		objects = append(objects, k8sObj)
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

// registerCustomSchemes registers additional API schemes
func registerCustomSchemes() error {
	// Register CRD scheme
	if err := apiextensionsv1.AddToScheme(scheme.Scheme); err != nil {
		return err
	}
	return nil
}

// ProcessObjects demonstrates how to work with different Kubernetes object types
func ProcessObjects(objects []K8sObject) {
	fmt.Printf("\n🎯 Detailed Object Analysis:\n")

	for i, obj := range objects {
		fmt.Printf("\n%d. %s: %s\n", i+1, obj.Kind, obj.Name)
		fmt.Printf("   File: %s\n", filepath.Base(obj.File))

		// Handle different types with type switch
		switch typed := obj.Object.(type) {
		// Apps
		case *appsv1.Deployment:
			fixtures.AddDeployment(typed)
		case *appsv1.StatefulSet:
			fixtures.AddStatefulSet(typed)
		case *appsv1.DaemonSet:
			fixtures.AddDaemonSet(typed)
		case *appsv1.ReplicaSet:
			fixtures.AddReplicaSet(typed)

		// Core
		case *corev1.Pod:
			fixtures.AddPod(typed)
		case *corev1.Service:
			fixtures.AddService(typed)
		case *corev1.ConfigMap:
			fixtures.AddConfigMap(typed)
		case *corev1.Secret:
			fixtures.AddSecret(typed)
		case *corev1.Namespace:
			fixtures.AddNamespace(typed)
		case *corev1.ServiceAccount:
			fixtures.AddServiceAccount(typed)
		case *corev1.PersistentVolume:
			fixtures.AddPersistentVolume(typed)
		case *corev1.PersistentVolumeClaim:
			fixtures.AddPersistentVolumeClaim(typed)

			// Batch
		case *batchv1.Job:
			fixtures.AddJob(typed)
		case *batchv1.CronJob:
			fixtures.AddCronJob(typed)

			// Networking
		case *networkingv1.Ingress:
			fixtures.AddIngress(typed)
		case *networkingv1.NetworkPolicy:
			fixtures.AddNetworkPolicy(typed)
		// RBAC
		case *rbacv1.Role:
			fixtures.AddRole(typed)
		case *rbacv1.RoleBinding:
			fixtures.AddRoleBinding(typed)
		case *rbacv1.ClusterRole:
			fixtures.AddClusterRole(typed)
		case *rbacv1.ClusterRoleBinding:
			fixtures.AddClusterRoleBinding(typed)

		// Storage
		case *storagev1.StorageClass:
			fixtures.AddStorageClass(typed)

		// Extensions
		case *apiextensionsv1.CustomResourceDefinition:
			fixtures.AddCustomResourceDefinition(typed)

		default:
			handleUnknownType(obj.Kind, typed)
		}
	}
	// Print summary
	fmt.Printf("\n📊 Created fixtures - Total: %d\n", fixtures.GetTotalCount())
	fmt.Printf("  Deployments: %d\n", len(fixtures.GetDeployments()))
	fmt.Printf("  Services: %d\n", len(fixtures.GetServices()))
	fmt.Printf("  ServiceAccounts: %d\n", len(fixtures.GetServiceAccounts()))
	fmt.Printf("  Pods: %d\n", len(fixtures.GetPods()))
	fmt.Printf("  ConfigMaps: %d\n", len(fixtures.GetConfigMaps()))
	fmt.Printf("  Secrets: %d\n", len(fixtures.GetSecrets()))
	// Save fixtures to file
	err := fixtures.SaveFixtures()
	if err != nil {
		fmt.Printf("Failed to save fixtures: %v", err)
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

// Apps Handlers
func handleDeployment(d *appsv1.Deployment) {
	fmt.Printf("   📦 Deployment Details:\n")
	if d.Spec.Replicas != nil {
		fmt.Printf("      Replicas: %d\n", *d.Spec.Replicas)
	} else {
		fmt.Printf("      Replicas: not specified (defaults to 1)\n")
	}
	if d.Spec.Strategy.Type != "" {
		fmt.Printf("      Strategy: %s\n", d.Spec.Strategy.Type)
	}
	if len(d.Spec.Template.Spec.Containers) > 0 {
		fmt.Printf("      Containers: %d\n", len(d.Spec.Template.Spec.Containers))
		for i, c := range d.Spec.Template.Spec.Containers {
			fmt.Printf("        %d. %s (%s)\n", i+1, c.Name, c.Image)
		}
	}
}

func handleStatefulSet(ss *appsv1.StatefulSet) {
	fmt.Printf("   🗳️  StatefulSet Details:\n")
	if ss.Spec.Replicas != nil {
		fmt.Printf("      Replicas: %d\n", *ss.Spec.Replicas)
	}
	if ss.Spec.ServiceName != "" {
		fmt.Printf("      Service Name: %s\n", ss.Spec.ServiceName)
	}
}

func handleDaemonSet(ds *appsv1.DaemonSet) {
	fmt.Printf("   👹 DaemonSet Details:\n")
	if ds.Spec.UpdateStrategy.Type != "" {
		fmt.Printf("      Update Strategy: %s\n", ds.Spec.UpdateStrategy.Type)
	}
}

func handleReplicaSet(rs *appsv1.ReplicaSet) {
	fmt.Printf("   🔄 ReplicaSet Details:\n")
	if rs.Spec.Replicas != nil {
		fmt.Printf("      Replicas: %d\n", *rs.Spec.Replicas)
	}
	fmt.Printf("      Ready Replicas: %d\n", rs.Status.ReadyReplicas)
}

// Core Handlers
func handlePod(p *corev1.Pod) {
	fmt.Printf("   🚀 Pod Details:\n")
	fmt.Printf("      Phase: %s\n", p.Status.Phase)
	if len(p.Spec.Containers) > 0 {
		fmt.Printf("      Containers: %d\n", len(p.Spec.Containers))
		for i, c := range p.Spec.Containers {
			fmt.Printf("        %d. %s (%s)\n", i+1, c.Name, c.Image)
		}
	}
}

func handleService(s *corev1.Service) {
	fmt.Printf("   🌐 Service Details:\n")
	fmt.Printf("      Type: %s\n", s.Spec.Type)
	if len(s.Spec.Ports) > 0 {
		fmt.Printf("      Ports: %d\n", len(s.Spec.Ports))
		for i, p := range s.Spec.Ports {
			portInfo := fmt.Sprintf("%d. %s:%d", i+1, p.Name, p.Port)
			if p.TargetPort.IntVal != 0 {
				portInfo += fmt.Sprintf(" → %d", p.TargetPort.IntVal)
			} else if p.TargetPort.StrVal != "" {
				portInfo += fmt.Sprintf(" → %s", p.TargetPort.StrVal)
			}
			fmt.Printf("        %s\n", portInfo)
		}
	}
}

func handleConfigMap(cm *corev1.ConfigMap) {
	fmt.Printf("   📝 ConfigMap Details:\n")
	fmt.Printf("      Data Keys: %d\n", len(cm.Data))
	fmt.Printf("      Binary Data Keys: %d\n", len(cm.BinaryData))
	if len(cm.Data) > 0 {
		for key := range cm.Data {
			fmt.Printf("        - %s\n", key)
		}
	}
}

func handleSecret(s *corev1.Secret) {
	fmt.Printf("   🔐 Secret Details:\n")
	fmt.Printf("      Type: %s\n", s.Type)
	fmt.Printf("      Data Keys: %d\n", len(s.Data))
	fmt.Printf("      String Data Keys: %d\n", len(s.StringData))
}

func handleNamespace(ns *corev1.Namespace) {
	fmt.Printf("   📛 Namespace Details:\n")
	fmt.Printf("      Phase: %s\n", ns.Status.Phase)
}

func handleServiceAccount(sa *corev1.ServiceAccount) {
	fmt.Printf("   👤 ServiceAccount Details:\n")
	fmt.Printf("      Secrets: %d\n", len(sa.Secrets))
	if len(sa.Secrets) > 0 {
		for i, secret := range sa.Secrets {
			fmt.Printf("        %d. %s\n", i+1, secret.Name)
		}
	}
	if len(sa.ImagePullSecrets) > 0 {
		fmt.Printf("      Image Pull Secrets: %d\n", len(sa.ImagePullSecrets))
		for i, secret := range sa.ImagePullSecrets {
			fmt.Printf("        %d. %s\n", i+1, secret.Name)
		}
	}
	// Check for automountServiceAccountToken
	if sa.AutomountServiceAccountToken != nil {
		fmt.Printf("      Automount Token: %v\n", *sa.AutomountServiceAccountToken)
	}
}

func handlePersistentVolume(pv *corev1.PersistentVolume) {
	fmt.Printf("   💽 PersistentVolume Details:\n")
	fmt.Printf("      Status: %s\n", pv.Status.Phase)
	fmt.Printf("      Capacity: %v\n", pv.Spec.Capacity.Storage())
	fmt.Printf("      Access Modes: %v\n", pv.Spec.AccessModes)
	if pv.Spec.ClaimRef != nil {
		fmt.Printf("      Claim: %s/%s\n", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
	}
}

func handlePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) {
	fmt.Printf("   💾 PersistentVolumeClaim Details:\n")
	fmt.Printf("      Status: %s\n", pvc.Status.Phase)
	if pvc.Spec.Resources.Requests != nil {
		fmt.Printf("      Storage: %v\n", pvc.Spec.Resources.Requests.Storage())
	}
	fmt.Printf("      Access Modes: %v\n", pvc.Spec.AccessModes)
}

func handleResourceQuota(rq *corev1.ResourceQuota) {
	fmt.Printf("   📊 ResourceQuota Details:\n")
	fmt.Printf("      Hard Limits: %d\n", len(rq.Spec.Hard))
	for resource := range rq.Spec.Hard {
		fmt.Printf("        - %s: %v\n", resource, rq.Spec.Hard[resource])
	}
}

func handleLimitRange(lr *corev1.LimitRange) {
	fmt.Printf("   ⚖️  LimitRange Details:\n")
	fmt.Printf("      Limits: %d\n", len(lr.Spec.Limits))
	for i, limit := range lr.Spec.Limits {
		fmt.Printf("        %d. Type: %s\n", i+1, limit.Type)
	}
}

// Batch Handlers
func handleJob(j *batchv1.Job) {
	fmt.Printf("   ⚙️  Job Details:\n")
	if j.Spec.Completions != nil {
		fmt.Printf("      Completions: %d\n", *j.Spec.Completions)
	}
	if j.Spec.Parallelism != nil {
		fmt.Printf("      Parallelism: %d\n", *j.Spec.Parallelism)
	}
}

func handleCronJob(cj *batchv1.CronJob) {
	fmt.Printf("   ⏰ CronJob Details:\n")
	fmt.Printf("      Schedule: %s\n", cj.Spec.Schedule)
	if cj.Spec.StartingDeadlineSeconds != nil {
		fmt.Printf("      Starting Deadline Seconds: %d\n", *cj.Spec.StartingDeadlineSeconds)
	}
}

// Networking Handlers
func handleIngress(i *networkingv1.Ingress) {
	fmt.Printf("   🚪 Ingress Details:\n")
	fmt.Printf("      Rules: %d\n", len(i.Spec.Rules))
	if i.Spec.DefaultBackend != nil {
		fmt.Printf("      Default Backend: configured\n")
	}
	if i.Spec.TLS != nil {
		fmt.Printf("      TLS Hosts: %d\n", len(i.Spec.TLS))
	}
}

func handleNetworkPolicy(np *networkingv1.NetworkPolicy) {
	fmt.Printf("   🛡️  NetworkPolicy Details:\n")
	fmt.Printf("      Policy Types: %v\n", np.Spec.PolicyTypes)
	fmt.Printf("      Pod Selector: %v\n", np.Spec.PodSelector)
}

// RBAC Handlers
func handleRole(r *rbacv1.Role) {
	fmt.Printf("   👤 Role Details:\n")
	fmt.Printf("      Rules: %d\n", len(r.Rules))
	for i, rule := range r.Rules {
		fmt.Printf("        %d. Resources: %v, Verbs: %v\n", i+1, rule.Resources, rule.Verbs)
	}
}

func handleRoleBinding(rb *rbacv1.RoleBinding) {
	fmt.Printf("   🔗 RoleBinding Details:\n")
	fmt.Printf("      Subjects: %d\n", len(rb.Subjects))
	fmt.Printf("      Role Ref: %s (%s)\n", rb.RoleRef.Name, rb.RoleRef.Kind)
	for i, subject := range rb.Subjects {
		fmt.Printf("        %d. %s (%s)\n", i+1, subject.Name, subject.Kind)
	}
}

func handleClusterRole(cr *rbacv1.ClusterRole) {
	fmt.Printf("   👥 ClusterRole Details:\n")
	fmt.Printf("      Rules: %d\n", len(cr.Rules))
}

func handleClusterRoleBinding(crb *rbacv1.ClusterRoleBinding) {
	fmt.Printf("   🔗 ClusterRoleBinding Details:\n")
	fmt.Printf("      Subjects: %d\n", len(crb.Subjects))
	fmt.Printf("      Role Ref: %s (%s)\n", crb.RoleRef.Name, crb.RoleRef.Kind)
}

// Storage Handlers
func handleStorageClass(sc *storagev1.StorageClass) {
	fmt.Printf("   💿 StorageClass Details:\n")
	fmt.Printf("      Provisioner: %s\n", sc.Provisioner)
	fmt.Printf("      Parameters: %d\n", len(sc.Parameters))
	for key, value := range sc.Parameters {
		fmt.Printf("        - %s: %s\n", key, value)
	}
}

// Extensions Handlers
func handleCRD(crd *apiextensionsv1.CustomResourceDefinition) {
	fmt.Printf("   📋 CustomResourceDefinition Details:\n")
	fmt.Printf("      Group: %s\n", crd.Spec.Group)
	fmt.Printf("      Kind: %s\n", crd.Spec.Names.Kind)
	fmt.Printf("      Scope: %s\n", crd.Spec.Scope)
	fmt.Printf("      Versions: %d\n", len(crd.Spec.Versions))
}

func handleUnknownType(kind string, obj interface{}) {
	fmt.Printf("   ❓ Unknown Type: %s (%T)\n", kind, obj)
}
