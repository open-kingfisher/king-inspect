package basic

import (
	"sync"

	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&defaultNamespaceCheck{})
}

type defaultNamespaceCheck struct{}

type alert struct {
	diagnostics []check.Diagnostic
	mu          sync.Mutex
}

func (alert *alert) GetDiagnostics() []check.Diagnostic {
	return alert.diagnostics
}

func (alert *alert) SetDiagnostics(d []check.Diagnostic) {
	alert.diagnostics = d
}

func (alert *alert) warn(k8stype check.Kind, itemMeta metav1.ObjectMeta) {
	d := check.Diagnostic{
		Severity: check.Warning,
		Message:  check.Message[107],
		Kind:     k8stype,
		Object:   &itemMeta,
		Owners:   itemMeta.GetOwnerReferences(),
	}
	alert.mu.Lock()
	alert.diagnostics = append(alert.diagnostics, d)
	alert.mu.Unlock()
}

// Name 返回此检查的唯一名称
func (nc *defaultNamespaceCheck) Name() string {
	return "default-namespace"
}

// Groups 返回此检查应属于的组名列表
func (nc *defaultNamespaceCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (nc *defaultNamespaceCheck) Description() string {
	return "检查是否有用户在缺省名称空间中创建了k8s对象"
}

// checkPods check if there are pods in the default namespace
func (nc *defaultNamespaceCheck) checkPods(items *corev1.PodList, alert *alert) {
	for _, item := range items.Items {
		if corev1.NamespaceDefault == item.GetNamespace() {
			alert.warn(check.Pod, item.ObjectMeta)
		}
	}
}

// checkPodTemplates check if there are pod templates in the default namespace
func (nc *defaultNamespaceCheck) checkPodTemplates(items *corev1.PodTemplateList, alert *alert) {
	for _, item := range items.Items {
		if corev1.NamespaceDefault == item.GetNamespace() {
			alert.warn(check.PodTemplate, item.ObjectMeta)
		}
	}
}

// checkPVCs check if there are pvcs in the default namespace
func (nc *defaultNamespaceCheck) checkPVCs(items *corev1.PersistentVolumeClaimList, alert *alert) {
	for _, item := range items.Items {
		if corev1.NamespaceDefault == item.GetNamespace() {
			alert.warn(check.PersistentVolumeClaim, item.ObjectMeta)
		}
	}
}

// checkConfigMaps check if there are config maps in the default namespace
func (nc *defaultNamespaceCheck) checkConfigMaps(items *corev1.ConfigMapList, alert *alert) {
	for _, item := range items.Items {
		if corev1.NamespaceDefault == item.GetNamespace() {
			alert.warn(check.ConfigMap, item.ObjectMeta)
		}
	}
}

// checkervices check if there are user created services in the default namespace
func (nc *defaultNamespaceCheck) checkervices(items *corev1.ServiceList, alert *alert) {
	for _, item := range items.Items {
		if corev1.NamespaceDefault == item.GetNamespace() && item.GetName() != "kubernetes" {
			alert.warn(check.Service, item.ObjectMeta)
		}
	}
}

// checkecrets check if there are user created secrets in the default namespace
func (nc *defaultNamespaceCheck) checkecrets(items *corev1.SecretList, alert *alert) {
	for _, item := range items.Items {
		if corev1.NamespaceDefault == item.GetNamespace() && item.Type != corev1.SecretTypeServiceAccountToken {
			alert.warn(check.Secret, item.ObjectMeta)
		}
	}
}

// checkA check if there are user created SAs in the default namespace
func (nc *defaultNamespaceCheck) checkA(items *corev1.ServiceAccountList, alert *alert) {
	for _, item := range items.Items {
		if corev1.NamespaceDefault == item.GetNamespace() && item.GetName() != "default" {
			alert.warn(check.ServiceAccount, item.ObjectMeta)
		}
	}
}

// Run 运行这个检查
func (nc *defaultNamespaceCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	alert := &alert{}
	var g errgroup.Group
	g.Go(func() error {
		nc.checkPods(objects.Pods, alert)
		return nil
	})

	g.Go(func() error {
		nc.checkPodTemplates(objects.PodTemplates, alert)
		return nil
	})

	g.Go(func() error {
		nc.checkPVCs(objects.PersistentVolumeClaims, alert)
		return nil
	})

	g.Go(func() error {
		nc.checkConfigMaps(objects.ConfigMaps, alert)
		return nil
	})

	g.Go(func() error {
		nc.checkervices(objects.Services, alert)
		return nil
	})

	g.Go(func() error {
		nc.checkecrets(objects.Secrets, alert)
		return nil
	})

	g.Go(func() error {
		nc.checkA(objects.ServiceAccounts, alert)
		return nil
	})

	err := g.Wait()
	var summary check.Summary
	summary.Total = len(alert.GetDiagnostics())
	summary.Issue = summary.Total
	summary.Warning = summary.Total
	return alert.GetDiagnostics(), summary, err
}
