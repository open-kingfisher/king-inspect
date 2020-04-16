package basic

import (
	ar "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&webhookReplaementCheck{})
}

type webhookReplaementCheck struct{}

// Name 返回此检查的唯一名称
func (w *webhookReplaementCheck) Name() string {
	return "admission-controller-webhook-replacement"
}

// Groups 返回此检查应属于的组名列表
func (w *webhookReplaementCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (w *webhookReplaementCheck) Description() string {
	return "Check for admission control webhooks that could cause problems during upgrades or node replacement"
}

// Run runs this check on a set of checkrnetes objects.
func (w *webhookReplaementCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	const apiserverServiceName = "checkrnetes"

	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.ValidatingWebhookConfigurations.Items) + len(objects.MutatingWebhookConfigurations.Items)
	for _, config := range objects.ValidatingWebhookConfigurations.Items {
		for _, wh := range config.Webhooks {
			if *wh.FailurePolicy == ar.Ignore {
				// Webhooks with failurePolicy: Ignore are fine.
				continue
			}
			if wh.ClientConfig.Service == nil {
				// Webhooks whose targets are external to the cluster are fine.
				continue
			}
			if wh.ClientConfig.Service.Namespace == metav1.NamespaceDefault &&
				wh.ClientConfig.Service.Name == apiserverServiceName {
				// Webhooks that target the check-apiserver are fine.
				continue
			}
			if !selectorMatchesNamespace(wh.NamespaceSelector, objects.SystemNamespace) {
				// Webhooks that don't apply to check-system are fine.
				continue
			}
			var svcNamespace *v1.Namespace
			for _, ns := range objects.Namespaces.Items {
				if ns.Name == wh.ClientConfig.Service.Namespace {
					svcNamespace = &ns
				}
			}
			if svcNamespace != nil &&
				!selectorMatchesNamespace(wh.NamespaceSelector, svcNamespace) &&
				len(objects.Nodes.Items) > 1 {
				// Webhooks that don't apply to their own namespace are fine, as
				// long as there's more than one node in the cluster.
				continue
			}

			d := check.Diagnostic{
				Severity: check.Error,
				Message:  "Validating webhook is configured in such a way that it may be problematic during upgrades.",
				Kind:     check.ValidatingWebhookConfiguration,
				Object:   &config.ObjectMeta,
				Owners:   config.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)

			// We don't want to produce diagnostics for multiple webhooks in the
			// same webhook configuration, so break out of the inner loop if we
			// get here.
			break
		}
	}

	for _, config := range objects.MutatingWebhookConfigurations.Items {
		for _, wh := range config.Webhooks {
			if *wh.FailurePolicy == ar.Ignore {
				// Webhooks with failurePolicy: Ignore are fine.
				continue
			}
			if wh.ClientConfig.Service == nil {
				// Webhooks whose targets are external to the cluster are fine.
				continue
			}
			if wh.ClientConfig.Service.Namespace == metav1.NamespaceDefault &&
				wh.ClientConfig.Service.Name == apiserverServiceName {
				// Webhooks that target the check-apiserver are fine.
				continue
			}
			if !selectorMatchesNamespace(wh.NamespaceSelector, objects.SystemNamespace) {
				// Webhooks that don't apply to check-system are fine.
				continue
			}
			var svcNamespace *v1.Namespace
			for _, ns := range objects.Namespaces.Items {
				if ns.Name == wh.ClientConfig.Service.Namespace {
					svcNamespace = &ns
				}
			}
			if svcNamespace != nil &&
				!selectorMatchesNamespace(wh.NamespaceSelector, svcNamespace) &&
				len(objects.Nodes.Items) > 1 {
				// Webhooks that don't apply to their own namespace are fine, as
				// long as there's more than one node in the cluster.
				continue
			}

			d := check.Diagnostic{
				Severity: check.Error,
				Message:  "Mutating webhook is configured in such a way that it may be problematic during upgrades.",
				Kind:     check.MutatingWebhookConfiguration,
				Object:   &config.ObjectMeta,
				Owners:   config.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)

			// We don't want to produce diagnostics for multiple webhooks in the
			// same webhook configuration, so break out of the inner loop if we
			// get here.
			break
		}
	}
	summary.Issue = len(diagnostics)
	summary.Error = summary.Issue
	return diagnostics, summary, nil
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

func contains(list []string, name string) bool {
	for _, l := range list {
		if l == name {
			return true
		}
	}
	return false
}
