package fixtures

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/kubernetes/test/ctest/ctestglobals"
)

// All Kubernetes object types as separate slices
var (
	// Apps
	Deployments  []*appsv1.Deployment
	StatefulSets []*appsv1.StatefulSet
	DaemonSets   []*appsv1.DaemonSet
	ReplicaSets  []*appsv1.ReplicaSet

	// Core
	Pods                   []*corev1.Pod
	Services               []*corev1.Service
	ConfigMaps             []*corev1.ConfigMap
	Secrets                []*corev1.Secret
	Namespaces             []*corev1.Namespace
	ServiceAccounts        []*corev1.ServiceAccount
	PersistentVolumes      []*corev1.PersistentVolume
	PersistentVolumeClaims []*corev1.PersistentVolumeClaim
	ResourceQuotas         []*corev1.ResourceQuota
	LimitRanges            []*corev1.LimitRange

	// Batch
	Jobs     []*batchv1.Job
	CronJobs []*batchv1.CronJob

	// Networking
	Ingresses       []*networkingv1.Ingress
	NetworkPolicies []*networkingv1.NetworkPolicy

	// RBAC
	Roles               []*rbacv1.Role
	RoleBindings        []*rbacv1.RoleBinding
	ClusterRoles        []*rbacv1.ClusterRole
	ClusterRoleBindings []*rbacv1.ClusterRoleBinding

	// Storage
	StorageClasses []*storagev1.StorageClass

	// Extensions
	CustomResourceDefinitions []*apiextensionsv1.CustomResourceDefinition

	AllObjects []interface{}

	initialized bool

	mu sync.RWMutex
)

const fixturesFile = "./fixtures/" + ctestglobals.TestExternalFixtureFile

// InitializeFixtures clears all fixtures
func InitializeFixtures() {
	mu.Lock()
	defer mu.Unlock()

	// Clear all existing fixtures
	// Apps
	Deployments = []*appsv1.Deployment{}
	StatefulSets = []*appsv1.StatefulSet{}
	DaemonSets = []*appsv1.DaemonSet{}
	ReplicaSets = []*appsv1.ReplicaSet{}

	// Core
	Pods = []*corev1.Pod{}
	Services = []*corev1.Service{}
	ConfigMaps = []*corev1.ConfigMap{}
	Secrets = []*corev1.Secret{}
	Namespaces = []*corev1.Namespace{}
	ServiceAccounts = []*corev1.ServiceAccount{}
	PersistentVolumes = []*corev1.PersistentVolume{}
	PersistentVolumeClaims = []*corev1.PersistentVolumeClaim{}
	ResourceQuotas = []*corev1.ResourceQuota{}
	LimitRanges = []*corev1.LimitRange{}

	// Batch
	Jobs = []*batchv1.Job{}
	CronJobs = []*batchv1.CronJob{}

	// Networking
	Ingresses = []*networkingv1.Ingress{}
	NetworkPolicies = []*networkingv1.NetworkPolicy{}

	// RBAC
	Roles = []*rbacv1.Role{}
	RoleBindings = []*rbacv1.RoleBinding{}
	ClusterRoles = []*rbacv1.ClusterRole{}
	ClusterRoleBindings = []*rbacv1.ClusterRoleBinding{}

	// Storage
	StorageClasses = []*storagev1.StorageClass{}

	// Extensions
	CustomResourceDefinitions = []*apiextensionsv1.CustomResourceDefinition{}

	initialized = true
	fmt.Println("✅ Fixtures initialized!")
}

// ========== APPS ADDERS ==========
func AddDeployment(deployment *appsv1.Deployment) {
	mu.Lock()
	defer mu.Unlock()
	Deployments = append(Deployments, deployment)
	AllObjects = append(AllObjects, deployment)
}

func AddStatefulSet(statefulSet *appsv1.StatefulSet) {
	mu.Lock()
	defer mu.Unlock()
	StatefulSets = append(StatefulSets, statefulSet)
	AllObjects = append(AllObjects, statefulSet)
}

func AddDaemonSet(daemonSet *appsv1.DaemonSet) {
	mu.Lock()
	defer mu.Unlock()
	DaemonSets = append(DaemonSets, daemonSet)
	AllObjects = append(AllObjects, daemonSet)
}

func AddReplicaSet(replicaSet *appsv1.ReplicaSet) {
	mu.Lock()
	defer mu.Unlock()
	ReplicaSets = append(ReplicaSets, replicaSet)
	AllObjects = append(AllObjects, replicaSet)
}

// ========== CORE ADDERS ==========
func AddPod(pod *corev1.Pod) {
	mu.Lock()
	defer mu.Unlock()
	Pods = append(Pods, pod)
	AllObjects = append(AllObjects, pod)
}

func AddService(service *corev1.Service) {
	mu.Lock()
	defer mu.Unlock()
	Services = append(Services, service)
	AllObjects = append(AllObjects, service)
}

func AddConfigMap(configMap *corev1.ConfigMap) {
	mu.Lock()
	defer mu.Unlock()
	ConfigMaps = append(ConfigMaps, configMap)
	AllObjects = append(AllObjects, configMap)
}

func AddSecret(secret *corev1.Secret) {
	mu.Lock()
	defer mu.Unlock()
	Secrets = append(Secrets, secret)
	AllObjects = append(AllObjects, secret)
}

func AddNamespace(namespace *corev1.Namespace) {
	mu.Lock()
	defer mu.Unlock()
	Namespaces = append(Namespaces, namespace)
	AllObjects = append(AllObjects, namespace)
}

func AddServiceAccount(serviceAccount *corev1.ServiceAccount) {
	mu.Lock()
	defer mu.Unlock()
	ServiceAccounts = append(ServiceAccounts, serviceAccount)
	AllObjects = append(AllObjects, serviceAccount)
}

func AddPersistentVolume(pv *corev1.PersistentVolume) {
	mu.Lock()
	defer mu.Unlock()
	PersistentVolumes = append(PersistentVolumes, pv)
	AllObjects = append(AllObjects, pv)
}

func AddPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) {
	mu.Lock()
	defer mu.Unlock()
	PersistentVolumeClaims = append(PersistentVolumeClaims, pvc)
	AllObjects = append(AllObjects, pvc)
}

func AddResourceQuota(resourceQuota *corev1.ResourceQuota) {
	mu.Lock()
	defer mu.Unlock()
	ResourceQuotas = append(ResourceQuotas, resourceQuota)
	AllObjects = append(AllObjects, resourceQuota)
}

func AddLimitRange(limitRange *corev1.LimitRange) {
	mu.Lock()
	defer mu.Unlock()
	LimitRanges = append(LimitRanges, limitRange)
	AllObjects = append(AllObjects, limitRange)
}

// ========== BATCH ADDERS ==========
func AddJob(job *batchv1.Job) {
	mu.Lock()
	defer mu.Unlock()
	Jobs = append(Jobs, job)
	AllObjects = append(AllObjects, job)
}

func AddCronJob(cronJob *batchv1.CronJob) {
	mu.Lock()
	defer mu.Unlock()
	CronJobs = append(CronJobs, cronJob)
	AllObjects = append(AllObjects, cronJob)
}

// ========== NETWORKING ADDERS ==========
func AddIngress(ingress *networkingv1.Ingress) {
	mu.Lock()
	defer mu.Unlock()
	Ingresses = append(Ingresses, ingress)
	AllObjects = append(AllObjects, ingress)
}

func AddNetworkPolicy(networkPolicy *networkingv1.NetworkPolicy) {
	mu.Lock()
	defer mu.Unlock()
	NetworkPolicies = append(NetworkPolicies, networkPolicy)
	AllObjects = append(AllObjects, networkPolicy)
}

// ========== RBAC ADDERS ==========
func AddRole(role *rbacv1.Role) {
	mu.Lock()
	defer mu.Unlock()
	Roles = append(Roles, role)
	AllObjects = append(AllObjects, role)
}

func AddRoleBinding(roleBinding *rbacv1.RoleBinding) {
	mu.Lock()
	defer mu.Unlock()
	RoleBindings = append(RoleBindings, roleBinding)
	AllObjects = append(AllObjects, roleBinding)
}

func AddClusterRole(clusterRole *rbacv1.ClusterRole) {
	mu.Lock()
	defer mu.Unlock()
	ClusterRoles = append(ClusterRoles, clusterRole)
	AllObjects = append(AllObjects, clusterRole)
}

func AddClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) {
	mu.Lock()
	defer mu.Unlock()
	ClusterRoleBindings = append(ClusterRoleBindings, clusterRoleBinding)
	AllObjects = append(AllObjects, clusterRoleBinding)
}

// ========== STORAGE ADDERS ==========
func AddStorageClass(storageClass *storagev1.StorageClass) {
	mu.Lock()
	defer mu.Unlock()
	StorageClasses = append(StorageClasses, storageClass)
	AllObjects = append(AllObjects, storageClass)
}

// ========== EXTENSIONS ADDERS ==========
func AddCustomResourceDefinition(crd *apiextensionsv1.CustomResourceDefinition) {
	mu.Lock()
	defer mu.Unlock()
	CustomResourceDefinitions = append(CustomResourceDefinitions, crd)
	AllObjects = append(AllObjects, crd)
}

// ========== GETTERS ==========

// Apps Getters
func GetDeployments() []*appsv1.Deployment {
	mu.RLock()
	defer mu.RUnlock()
	return Deployments
}

func GetStatefulSets() []*appsv1.StatefulSet {
	mu.RLock()
	defer mu.RUnlock()
	return StatefulSets
}

func GetDaemonSets() []*appsv1.DaemonSet {
	mu.RLock()
	defer mu.RUnlock()
	return DaemonSets
}

func GetReplicaSets() []*appsv1.ReplicaSet {
	mu.RLock()
	defer mu.RUnlock()
	return ReplicaSets
}

// Core Getters
func GetPods() []*corev1.Pod {
	mu.RLock()
	defer mu.RUnlock()
	return Pods
}

func GetServices() []*corev1.Service {
	mu.RLock()
	defer mu.RUnlock()
	return Services
}

func GetConfigMaps() []*corev1.ConfigMap {
	mu.RLock()
	defer mu.RUnlock()
	return ConfigMaps
}

func GetSecrets() []*corev1.Secret {
	mu.RLock()
	defer mu.RUnlock()
	return Secrets
}

func GetNamespaces() []*corev1.Namespace {
	mu.RLock()
	defer mu.RUnlock()
	return Namespaces
}

func GetServiceAccounts() []*corev1.ServiceAccount {
	mu.RLock()
	defer mu.RUnlock()
	return ServiceAccounts
}

func GetPersistentVolumes() []*corev1.PersistentVolume {
	mu.RLock()
	defer mu.RUnlock()
	return PersistentVolumes
}

func GetPersistentVolumeClaims() []*corev1.PersistentVolumeClaim {
	mu.RLock()
	defer mu.RUnlock()
	return PersistentVolumeClaims
}

func GetResourceQuotas() []*corev1.ResourceQuota {
	mu.RLock()
	defer mu.RUnlock()
	return ResourceQuotas
}

func GetLimitRanges() []*corev1.LimitRange {
	mu.RLock()
	defer mu.RUnlock()
	return LimitRanges
}

// Batch Getters
func GetJobs() []*batchv1.Job {
	mu.RLock()
	defer mu.RUnlock()
	return Jobs
}

func GetCronJobs() []*batchv1.CronJob {
	mu.RLock()
	defer mu.RUnlock()
	return CronJobs
}

// Networking Getters
func GetIngresses() []*networkingv1.Ingress {
	mu.RLock()
	defer mu.RUnlock()
	return Ingresses
}

func GetNetworkPolicies() []*networkingv1.NetworkPolicy {
	mu.RLock()
	defer mu.RUnlock()
	return NetworkPolicies
}

// RBAC Getters
func GetRoles() []*rbacv1.Role {
	mu.RLock()
	defer mu.RUnlock()
	return Roles
}

func GetRoleBindings() []*rbacv1.RoleBinding {
	mu.RLock()
	defer mu.RUnlock()
	return RoleBindings
}

func GetClusterRoles() []*rbacv1.ClusterRole {
	mu.RLock()
	defer mu.RUnlock()
	return ClusterRoles
}

func GetClusterRoleBindings() []*rbacv1.ClusterRoleBinding {
	mu.RLock()
	defer mu.RUnlock()
	return ClusterRoleBindings
}

// Storage Getters
func GetStorageClasses() []*storagev1.StorageClass {
	mu.RLock()
	defer mu.RUnlock()
	return StorageClasses
}

// Extensions Getters
func GetCustomResourceDefinitions() []*apiextensionsv1.CustomResourceDefinition {
	mu.RLock()
	defer mu.RUnlock()
	return CustomResourceDefinitions
}

// Utility Getters
func GetAllObjects() []interface{} {
	mu.RLock()
	defer mu.RUnlock()
	return AllObjects
}

// func GetTotalCount() int {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	return len(AllObjects)
// }

// // GetCounts returns a map with counts of each type
// func GetCounts() map[string]int {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	return map[string]int{
// 		"Deployments":               len(Deployments),
// 		"StatefulSets":              len(StatefulSets),
// 		"DaemonSets":                len(DaemonSets),
// 		"ReplicaSets":               len(ReplicaSets),
// 		"Pods":                      len(Pods),
// 		"Services":                  len(Services),
// 		"ConfigMaps":                len(ConfigMaps),
// 		"Secrets":                   len(Secrets),
// 		"Namespaces":                len(Namespaces),
// 		"ServiceAccounts":           len(ServiceAccounts),
// 		"PersistentVolumes":         len(PersistentVolumes),
// 		"PersistentVolumeClaims":    len(PersistentVolumeClaims),
// 		"ResourceQuotas":            len(ResourceQuotas),
// 		"LimitRanges":               len(LimitRanges),
// 		"Jobs":                      len(Jobs),
// 		"CronJobs":                  len(CronJobs),
// 		"Ingresses":                 len(Ingresses),
// 		"NetworkPolicies":           len(NetworkPolicies),
// 		"Roles":                     len(Roles),
// 		"RoleBindings":              len(RoleBindings),
// 		"ClusterRoles":              len(ClusterRoles),
// 		"ClusterRoleBindings":       len(ClusterRoleBindings),
// 		"StorageClasses":            len(StorageClasses),
// 		"CustomResourceDefinitions": len(CustomResourceDefinitions),
// 	}
// }

func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return initialized
}

// // Add this function to check if we have any data
// func AreFixturesLoaded() bool {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	return len(Deployments) > 0 || len(Services) > 0 || len(ServiceAccounts) > 0
// }

// SaveFixtures saves all fixtures to a file
func SaveFixtures() error {
	mu.Lock()
	defer mu.Unlock()

	fixturesData := struct {
		Deployments               []*appsv1.Deployment                        `json:"deployments"`
		StatefulSets              []*appsv1.StatefulSet                       `json:"statefulSets"`
		DaemonSets                []*appsv1.DaemonSet                         `json:"daemonSets"`
		ReplicaSets               []*appsv1.ReplicaSet                        `json:"replicaSets"`
		Pods                      []*corev1.Pod                               `json:"pods"`
		Services                  []*corev1.Service                           `json:"services"`
		ConfigMaps                []*corev1.ConfigMap                         `json:"configMaps"`
		Secrets                   []*corev1.Secret                            `json:"secrets"`
		Namespaces                []*corev1.Namespace                         `json:"namespaces"`
		ServiceAccounts           []*corev1.ServiceAccount                    `json:"serviceAccounts"`
		PersistentVolumes         []*corev1.PersistentVolume                  `json:"persistentVolumes"`
		PersistentVolumeClaims    []*corev1.PersistentVolumeClaim             `json:"persistentVolumeClaims"`
		ResourceQuotas            []*corev1.ResourceQuota                     `json:"resourceQuotas"`
		LimitRanges               []*corev1.LimitRange                        `json:"limitRanges"`
		Jobs                      []*batchv1.Job                              `json:"jobs"`
		CronJobs                  []*batchv1.CronJob                          `json:"cronJobs"`
		Ingresses                 []*networkingv1.Ingress                     `json:"ingresses"`
		NetworkPolicies           []*networkingv1.NetworkPolicy               `json:"networkPolicies"`
		Roles                     []*rbacv1.Role                              `json:"roles"`
		RoleBindings              []*rbacv1.RoleBinding                       `json:"roleBindings"`
		ClusterRoles              []*rbacv1.ClusterRole                       `json:"clusterRoles"`
		ClusterRoleBindings       []*rbacv1.ClusterRoleBinding                `json:"clusterRoleBindings"`
		StorageClasses            []*storagev1.StorageClass                   `json:"storageClasses"`
		CustomResourceDefinitions []*apiextensionsv1.CustomResourceDefinition `json:"customResourceDefinitions"`
	}{
		Deployments:               Deployments,
		StatefulSets:              StatefulSets,
		DaemonSets:                DaemonSets,
		ReplicaSets:               ReplicaSets,
		Pods:                      Pods,
		Services:                  Services,
		ConfigMaps:                ConfigMaps,
		Secrets:                   Secrets,
		Namespaces:                Namespaces,
		ServiceAccounts:           ServiceAccounts,
		PersistentVolumes:         PersistentVolumes,
		PersistentVolumeClaims:    PersistentVolumeClaims,
		ResourceQuotas:            ResourceQuotas,
		LimitRanges:               LimitRanges,
		Jobs:                      Jobs,
		CronJobs:                  CronJobs,
		Ingresses:                 Ingresses,
		NetworkPolicies:           NetworkPolicies,
		Roles:                     Roles,
		RoleBindings:              RoleBindings,
		ClusterRoles:              ClusterRoles,
		ClusterRoleBindings:       ClusterRoleBindings,
		StorageClasses:            StorageClasses,
		CustomResourceDefinitions: CustomResourceDefinitions,
	}

	data, err := json.MarshalIndent(fixturesData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal fixtures: %w", err)
	}

	err = os.WriteFile(fixturesFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write fixtures file: %w", err)
	}

	totalCount := getTotalCount()
	fmt.Printf("✅ Fixtures saved to %s (%d objects)\n", fixturesFile, totalCount)
	return nil
}

// LoadFixtures loads fixtures from file
func LoadFixtures() error {
	mu.Lock()
	defer mu.Unlock()

	if _, err := os.Stat(fixturesFile); os.IsNotExist(err) {
		return fmt.Errorf("fixtures file not found: %s", fixturesFile)
	}

	data, err := os.ReadFile(fixturesFile)
	if err != nil {
		return fmt.Errorf("failed to read fixtures file: %w", err)
	}

	var fixturesData struct {
		Deployments               []*appsv1.Deployment                        `json:"deployments"`
		StatefulSets              []*appsv1.StatefulSet                       `json:"statefulSets"`
		DaemonSets                []*appsv1.DaemonSet                         `json:"daemonSets"`
		ReplicaSets               []*appsv1.ReplicaSet                        `json:"replicaSets"`
		Pods                      []*corev1.Pod                               `json:"pods"`
		Services                  []*corev1.Service                           `json:"services"`
		ConfigMaps                []*corev1.ConfigMap                         `json:"configMaps"`
		Secrets                   []*corev1.Secret                            `json:"secrets"`
		Namespaces                []*corev1.Namespace                         `json:"namespaces"`
		ServiceAccounts           []*corev1.ServiceAccount                    `json:"serviceAccounts"`
		PersistentVolumes         []*corev1.PersistentVolume                  `json:"persistentVolumes"`
		PersistentVolumeClaims    []*corev1.PersistentVolumeClaim             `json:"persistentVolumeClaims"`
		ResourceQuotas            []*corev1.ResourceQuota                     `json:"resourceQuotas"`
		LimitRanges               []*corev1.LimitRange                        `json:"limitRanges"`
		Jobs                      []*batchv1.Job                              `json:"jobs"`
		CronJobs                  []*batchv1.CronJob                          `json:"cronJobs"`
		Ingresses                 []*networkingv1.Ingress                     `json:"ingresses"`
		NetworkPolicies           []*networkingv1.NetworkPolicy               `json:"networkPolicies"`
		Roles                     []*rbacv1.Role                              `json:"roles"`
		RoleBindings              []*rbacv1.RoleBinding                       `json:"roleBindings"`
		ClusterRoles              []*rbacv1.ClusterRole                       `json:"clusterRoles"`
		ClusterRoleBindings       []*rbacv1.ClusterRoleBinding                `json:"clusterRoleBindings"`
		StorageClasses            []*storagev1.StorageClass                   `json:"storageClasses"`
		CustomResourceDefinitions []*apiextensionsv1.CustomResourceDefinition `json:"customResourceDefinitions"`
	}

	err = json.Unmarshal(data, &fixturesData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal fixtures: %w", err)
	}

	Deployments = fixturesData.Deployments
	StatefulSets = fixturesData.StatefulSets
	DaemonSets = fixturesData.DaemonSets
	ReplicaSets = fixturesData.ReplicaSets
	Pods = fixturesData.Pods
	Services = fixturesData.Services
	ConfigMaps = fixturesData.ConfigMaps
	Secrets = fixturesData.Secrets
	Namespaces = fixturesData.Namespaces
	ServiceAccounts = fixturesData.ServiceAccounts
	PersistentVolumes = fixturesData.PersistentVolumes
	PersistentVolumeClaims = fixturesData.PersistentVolumeClaims
	ResourceQuotas = fixturesData.ResourceQuotas
	LimitRanges = fixturesData.LimitRanges
	Jobs = fixturesData.Jobs
	CronJobs = fixturesData.CronJobs
	Ingresses = fixturesData.Ingresses
	NetworkPolicies = fixturesData.NetworkPolicies
	Roles = fixturesData.Roles
	RoleBindings = fixturesData.RoleBindings
	ClusterRoles = fixturesData.ClusterRoles
	ClusterRoleBindings = fixturesData.ClusterRoleBindings
	StorageClasses = fixturesData.StorageClasses
	CustomResourceDefinitions = fixturesData.CustomResourceDefinitions

	totalCount := getTotalCount()
	fmt.Printf("✅ Fixtures loaded from %s (%d objects)\n", fixturesFile, totalCount)
	return nil
}

// FixturesExist checks if the fixtures file exists
func FixturesExist() bool {
	_, err := os.Stat(fixturesFile)
	return err == nil
}

// ClearFixtures removes the fixtures file and clears memory
func ClearFixtures() error {
	mu.Lock()
	defer mu.Unlock()

	// Clear memory
	Deployments = nil
	StatefulSets = nil
	DaemonSets = nil
	ReplicaSets = nil
	Pods = nil
	Services = nil
	ConfigMaps = nil
	Secrets = nil
	Namespaces = nil
	ServiceAccounts = nil
	PersistentVolumes = nil
	PersistentVolumeClaims = nil
	ResourceQuotas = nil
	LimitRanges = nil
	Jobs = nil
	CronJobs = nil
	Ingresses = nil
	NetworkPolicies = nil
	Roles = nil
	RoleBindings = nil
	ClusterRoles = nil
	ClusterRoleBindings = nil
	StorageClasses = nil
	CustomResourceDefinitions = nil

	// Remove file
	if _, err := os.Stat(fixturesFile); err == nil {
		err := os.Remove(fixturesFile)
		if err != nil {
			return fmt.Errorf("failed to remove fixtures file: %w", err)
		}
		fmt.Println("✅ Fixtures cleared")
	}
	return nil
}

// AreFixturesLoaded checks if we have fixtures in memory
func AreFixturesLoaded() bool {
	mu.RLock()
	defer mu.RUnlock()
	return getTotalCount() > 0
}

// ========== UTILITY FUNCTIONS ==========

func GetCounts() map[string]int {
	mu.RLock()
	defer mu.RUnlock()
	return map[string]int{
		"Deployments":               len(Deployments),
		"StatefulSets":              len(StatefulSets),
		"DaemonSets":                len(DaemonSets),
		"ReplicaSets":               len(ReplicaSets),
		"Pods":                      len(Pods),
		"Services":                  len(Services),
		"ConfigMaps":                len(ConfigMaps),
		"Secrets":                   len(Secrets),
		"Namespaces":                len(Namespaces),
		"ServiceAccounts":           len(ServiceAccounts),
		"PersistentVolumes":         len(PersistentVolumes),
		"PersistentVolumeClaims":    len(PersistentVolumeClaims),
		"ResourceQuotas":            len(ResourceQuotas),
		"LimitRanges":               len(LimitRanges),
		"Jobs":                      len(Jobs),
		"CronJobs":                  len(CronJobs),
		"Ingresses":                 len(Ingresses),
		"NetworkPolicies":           len(NetworkPolicies),
		"Roles":                     len(Roles),
		"RoleBindings":              len(RoleBindings),
		"ClusterRoles":              len(ClusterRoles),
		"ClusterRoleBindings":       len(ClusterRoleBindings),
		"StorageClasses":            len(StorageClasses),
		"CustomResourceDefinitions": len(CustomResourceDefinitions),
	}
}

func GetTotalCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return getTotalCount()
}

// Helper function to calculate total count
func getTotalCount() int {
	return len(Deployments) + len(StatefulSets) + len(DaemonSets) + len(ReplicaSets) +
		len(Pods) + len(Services) + len(ConfigMaps) + len(Secrets) + len(Namespaces) +
		len(ServiceAccounts) + len(PersistentVolumes) + len(PersistentVolumeClaims) +
		len(ResourceQuotas) + len(LimitRanges) + len(Jobs) + len(CronJobs) +
		len(Ingresses) + len(NetworkPolicies) + len(Roles) + len(RoleBindings) +
		len(ClusterRoles) + len(ClusterRoleBindings) + len(StorageClasses) +
		len(CustomResourceDefinitions)
}
