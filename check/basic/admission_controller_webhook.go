package basic

import (
	"github.com/open-kingfisher/king-inspect/check"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"
)

func init() {
	check.Register(&webhookCheck{})
}

type webhookCheck struct{}

// Name 返回此检查的唯一名称
func (w *webhookCheck) Name() string {
	return "admission-controller-webhook"
}

// Groups 返回此检查应属于的组名列表
func (w *webhookCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (w *webhookCheck) Description() string {
	return "检查准入控制中的Webhook"
}

// Run 运行这个检查
func (w *webhookCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	const apiserverServiceName = "kubernetes"

	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.ValidatingWebhookConfigurations.Items) + len(objects.MutatingWebhookConfigurations.Items)
	for _, config := range objects.ValidatingWebhookConfigurations.Items {
		for _, wh := range config.Webhooks {
			if wh.ClientConfig.Service != nil {
				if !namespaceExists(objects.Namespaces, wh.ClientConfig.Service.Namespace) {
					diagnostics = append(diagnostics, check.Diagnostic{
						Severity: check.Error,
						Message:  check.Message[115],
						Kind:     check.ValidatingWebhookConfiguration,
						Object:   &config.ObjectMeta,
						Owners:   config.ObjectMeta.GetOwnerReferences(),
					})
					continue
				}

				if !serviceExists(objects.Services, wh.ClientConfig.Service.Name, wh.ClientConfig.Service.Namespace) {
					diagnostics = append(diagnostics, check.Diagnostic{
						Severity: check.Error,
						Message:  check.Message[116],
						Kind:     check.ValidatingWebhookConfiguration,
						Object:   &config.ObjectMeta,
						Owners:   config.ObjectMeta.GetOwnerReferences(),
					})
				}
			}
			if wh.NamespaceSelector != nil {
				if selectorMatchesNamespace(wh.NamespaceSelector, objects.SystemNamespace) {
					// Webhooks不应该应该应用到kube-system命名空间中
					diagnostics = append(diagnostics, check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[121],
						Kind:     check.ValidatingWebhookConfiguration,
						Object:   &config.ObjectMeta,
						Owners:   config.ObjectMeta.GetOwnerReferences(),
					})
				}
			}
		}
	}

	for _, config := range objects.MutatingWebhookConfigurations.Items {
		for _, wh := range config.Webhooks {
			if wh.ClientConfig.Service != nil {
				// Ensure that the service (and its namespace) that is configure actually exists.

				if !namespaceExists(objects.Namespaces, wh.ClientConfig.Service.Namespace) {
					diagnostics = append(diagnostics, check.Diagnostic{
						Severity: check.Error,
						Message:  check.Message[117],
						Kind:     check.MutatingWebhookConfiguration,
						Object:   &config.ObjectMeta,
						Owners:   config.ObjectMeta.GetOwnerReferences(),
					})
					continue
				}

				if !serviceExists(objects.Services, wh.ClientConfig.Service.Name, wh.ClientConfig.Service.Namespace) {
					diagnostics = append(diagnostics, check.Diagnostic{
						Severity: check.Error,
						Message:  check.Message[118],
						Kind:     check.MutatingWebhookConfiguration,
						Object:   &config.ObjectMeta,
						Owners:   config.ObjectMeta.GetOwnerReferences(),
					})
				}
			}
			if wh.NamespaceSelector != nil {
				if selectorMatchesNamespace(wh.NamespaceSelector, objects.SystemNamespace) {
					// Webhooks不应该应该应用到kube-system命名空间中
					diagnostics = append(diagnostics, check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[122],
						Kind:     check.MutatingWebhookConfiguration,
						Object:   &config.ObjectMeta,
						Owners:   config.ObjectMeta.GetOwnerReferences(),
					})
				}
			}
		}
	}
	summary.Issue = len(diagnostics)
	summary.Error = summary.Issue
	return diagnostics, summary, nil
}

func namespaceExists(namespaceList *v1.NamespaceList, namespace string) bool {
	for _, ns := range namespaceList.Items {
		if ns.Name == namespace {
			return true
		}
	}
	return false
}

func serviceExists(serviceList *v1.ServiceList, service, namespace string) bool {
	for _, svc := range serviceList.Items {
		if svc.Name == service && svc.Namespace == namespace {
			return true
		}
	}
	return false
}

func contains(list []string, name string) bool {
	for _, l := range list {
		if l == name {
			return true
		}
	}
	return false
}

func match(labels map[string]string, lbr metav1.LabelSelectorRequirement) bool {
	switch lbr.Operator {
	case metav1.LabelSelectorOpExists:
		if _, ok := labels[lbr.Key]; ok {
			return true
		}
		return false
	case metav1.LabelSelectorOpDoesNotExist:
		if _, ok := labels[lbr.Key]; !ok {
			return true
		}
		return false
	case metav1.LabelSelectorOpIn:
		if v, ok := labels[lbr.Key]; ok && contains(lbr.Values, v) {
			return true
		}
		return false
	case metav1.LabelSelectorOpNotIn:
		if v, ok := labels[lbr.Key]; !ok || !contains(lbr.Values, v) {
			return true
		}
		return false
	}
	return false
}

func selectorMatchesNamespace(selector *metav1.LabelSelector, namespace *corev1.Namespace) bool {
	if selector.Size() == 0 {
		return true
	}
	labels := namespace.GetLabels()
	for key, value := range selector.MatchLabels {
		if v, ok := labels[key]; !ok || v != value {
			return false
		}
	}
	for _, lbr := range selector.MatchExpressions {
		if !match(labels, lbr) {
			return false
		}
	}
	return true
}
