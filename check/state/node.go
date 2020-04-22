package state

import (
	"github.com/open-kingfisher/king-inspect/check"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	check.Register(&nodeCheck{})
}

type nodeCheck struct{}

// Name 返回此检查的唯一名称
func (c *nodeCheck) Name() string {
	return "node-state"
}

// Groups 返回此检查应属于的组名列表
func (c *nodeCheck) Groups() []string {
	return []string{"state"}
}

// Description 返回此检查的描述信息
func (c *nodeCheck) Description() string {
	return "检查集群中节点的状态"
}

// Run 运行这个检查
func (c *nodeCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Nodes.Items)
	for _, node := range objects.Nodes.Items {
		node := node
		if check.IsEnabled(c.Name(), &node.ObjectMeta) {
			for _, n := range node.Status.Conditions {
				if n.Status == corev1.ConditionUnknown {
					summary.Warning += 1
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[300],
						Kind:     check.Node,
						Object:   &node.ObjectMeta,
						Owners:   node.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
				}
				if n.Type == corev1.NodeReady && n.Status == corev1.ConditionFalse {
					summary.Error += 1
					d := check.Diagnostic{
						Severity: check.Error,
						Message:  check.Message[301],
						Kind:     check.Node,
						Object:   &node.ObjectMeta,
						Owners:   node.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
				}
				if n.Type == corev1.NodeMemoryPressure && n.Status == corev1.ConditionTrue {
					summary.Warning += 1
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[302],
						Kind:     check.Node,
						Object:   &node.ObjectMeta,
						Owners:   node.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
				}
				if n.Type == corev1.NodeDiskPressure && n.Status == corev1.ConditionTrue {
					summary.Warning += 1
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[303],
						Kind:     check.Node,
						Object:   &node.ObjectMeta,
						Owners:   node.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
				}
				if n.Type == corev1.NodePIDPressure && n.Status == corev1.ConditionTrue {
					summary.Warning += 1
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[304],
						Kind:     check.Node,
						Object:   &node.ObjectMeta,
						Owners:   node.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
				}
				if n.Type == corev1.NodeNetworkUnavailable && n.Status == corev1.ConditionTrue {
					summary.Warning += 1
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[305],
						Kind:     check.Node,
						Object:   &node.ObjectMeta,
						Owners:   node.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
				}
			}
		}
	}
	summary.Issue = len(diagnostics)
	return diagnostics, summary, nil
}
