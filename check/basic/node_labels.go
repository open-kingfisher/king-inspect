package basic

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&nodeLabelsTaintsCheck{})
}

type nodeLabelsTaintsCheck struct{}

// Name 返回此检查的唯一名称
func (*nodeLabelsTaintsCheck) Name() string {
	return "node-labels"
}

// Groups 返回此检查应属于的组名列表
func (*nodeLabelsTaintsCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (*nodeLabelsTaintsCheck) Description() string {
	return "检查节点是否有自定义标签"
}

// Run 运行这个检查
func (c *nodeLabelsTaintsCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Nodes.Items)
	for _, node := range objects.Nodes.Items {
		node := node
		if check.IsEnabled(c.Name(), &node.ObjectMeta) {
			for labelKey := range node.Labels {
				if !ischeckrnetesLabel(labelKey) {
					d := check.Diagnostic{
						Severity: check.Suggestion,
						Message:  check.Message[108],
						Kind:     check.Node,
						Object:   &node.ObjectMeta,
					}
					diagnostics = append(diagnostics, d)
					// Produce only one label diagnostic per node.
					break
				}
			}
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func ischeckrnetesLabel(key string) bool {
	// Built-in checkrnetes labels are in various subdomains of
	// checkrnetes.io. Assume all such labels are built in.
	return strings.Contains(key, corev1.ResourceDefaultNamespacePrefix)
}
