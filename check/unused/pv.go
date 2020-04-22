package unused

import (
	"github.com/open-kingfisher/king-inspect/check"
)

func init() {
	check.Register(&unusedPVCheck{})
}

type unusedPVCheck struct{}

// Name 返回此检查的唯一名称
func (p *unusedPVCheck) Name() string {
	return "unused-pv"
}

// Groups 返回此检查应属于的组名列表
func (p *unusedPVCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (p *unusedPVCheck) Description() string {
	return "检查集群中未使用的PV"
}

// Run 运行这个检查
func (p *unusedPVCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.PersistentVolumes.Items)
	for _, pv := range objects.PersistentVolumes.Items {
		pv := pv
		if pv.Spec.ClaimRef == nil && check.IsEnabled(p.Name(), &pv.ObjectMeta) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[200],
				Kind:     check.PersistentVolume,
				Object:   &pv.ObjectMeta,
				Owners:   pv.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
