package basic

import (
	"github.com/open-kingfisher/king-inspect/check"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	check.Register(&podSelectorCheck{})
}

type podSelectorCheck struct{}

// Name 返回此检查的唯一名称
func (p *podSelectorCheck) Name() string {
	return "node-name-pod-selector"
}

// Groups 返回此检查应属于的组名列表
func (p *podSelectorCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (p *podSelectorCheck) Description() string {
	return "检查Pod是否有节点选择标签使用的是kubernetes.io/hostname标签"
}

// Run 运行这个检查
func (p *podSelectorCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		nodeSelectorMap := pod.Spec.NodeSelector
		if _, ok := nodeSelectorMap[corev1.LabelHostname]; ok && check.IsEnabled(p.Name(), &pod.ObjectMeta) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[119],
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
