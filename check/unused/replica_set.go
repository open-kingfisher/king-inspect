package unused

import (
	"github.com/open-kingfisher/king-inspect/check"
)

func init() {
	check.Register(&unusedReplicaCheck{})
}

type unusedReplicaCheck struct{}

// Name 返回此检查的唯一名称
func (c *unusedReplicaCheck) Name() string {
	return "unused-replica-set"
}

// Groups 返回此检查应属于的组名列表
func (c *unusedReplicaCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (c *unusedReplicaCheck) Description() string {
	return "检查集群中没用使用的副本集"
}

// Run 运行这个检查
func (c *unusedReplicaCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	used, err := check.DeploymentReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	statefulSetRefs, err := check.StatefulSetReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	daemonSetRefs, err := check.DaemonSetReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	for k, v := range statefulSetRefs {
		used[k] = v
	}
	for k, v := range daemonSetRefs {
		used[k] = v
	}
	var summary check.Summary
	summary.Total = len(objects.ReplicaSets.Items)
	for _, re := range objects.ReplicaSets.Items {
		if re.OwnerReferences == nil && check.IsEnabled(c.Name(), &re.ObjectMeta) {
			re := re
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[205],
				Kind:     check.ReplicaSet,
				Object:   &re.ObjectMeta,
				Owners:   re.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		} else {
			for _, owner := range re.OwnerReferences {
				if _, ok := used[check.Identifier{Name: owner.Name, Namespace: re.GetNamespace()}]; !ok && check.IsEnabled(c.Name(), &re.ObjectMeta) {
					re := re
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[205],
						Kind:     check.ReplicaSet,
						Object:   &re.ObjectMeta,
						Owners:   re.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
					break
				}
			}
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
