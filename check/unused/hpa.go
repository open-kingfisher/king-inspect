package unused

import (
	"github.com/open-kingfisher/king-inspect/check"
)

func init() {
	check.Register(&unusedHPACheck{})
}

type unusedHPACheck struct{}

// Name 返回此检查的唯一名称
func (c *unusedHPACheck) Name() string {
	return "unused-hpa"
}

// Groups 返回此检查应属于的组名列表
func (c *unusedHPACheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (c *unusedHPACheck) Description() string {
	return "检查集群中没用使用的HPA"
}

// Run 运行这个检查
func (c *unusedHPACheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	used, err := check.DeploymentReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	statefulSeRetfs, err := check.StatefulSetReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	for k, v := range statefulSeRetfs {
		used[k] = v
	}
	var summary check.Summary
	summary.Total = len(objects.HPA.Items)
	for _, hpa := range objects.HPA.Items {
		if _, ok := used[check.Identifier{Name: hpa.GetName(), Namespace: hpa.GetNamespace()}]; !ok && check.IsEnabled(c.Name(), &hpa.ObjectMeta) {
			hpa := hpa
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[204],
				Kind:     check.HPA,
				Object:   &hpa.ObjectMeta,
				Owners:   hpa.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
