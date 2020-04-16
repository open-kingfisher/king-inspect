package state

import (
	corev1 "k8s.io/api/core/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&namespaceCheck{})
}

type namespaceCheck struct{}

// Name 返回此检查的唯一名称
func (c *namespaceCheck) Name() string {
	return "namespace-state"
}

// Groups 返回此检查应属于的组名列表
func (c *namespaceCheck) Groups() []string {
	return []string{"state"}
}

// Description 返回此检查的描述信息
func (c *namespaceCheck) Description() string {
	return "检查集群命名空间的状态"
}

// Run 运行这个检查
func (c *namespaceCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Namespaces.Items)
	for _, namespace := range objects.Namespaces.Items {
		namespace := namespace
		if !isNSActive(namespace.Status.Phase) && check.IsEnabled(c.Name(), &namespace.ObjectMeta) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[306],
				Kind:     check.Node,
				Object:   &namespace.ObjectMeta,
				Owners:   namespace.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}

	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func isNSActive(phase corev1.NamespacePhase) bool {
	return phase == corev1.NamespaceActive
}
