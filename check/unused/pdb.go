package unused

import (
	"kingfisher/king-inspect/check"
	"strings"
)

func init() {
	check.Register(&unusedPDBCheck{})
}

type unusedPDBCheck struct{}

// Name 返回此检查的唯一名称
func (p *unusedPDBCheck) Name() string {
	return "unused-pod-disruption-budget"
}

// Groups 返回此检查应属于的组名列表
func (p *unusedPDBCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (p *unusedPDBCheck) Description() string {
	return "检查集群中未使用的Pod中断预算"
}

// Run 运行这个检查
func (p *unusedPDBCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary

	summary.Total = len(objects.PodDisruptionBudgets.Items)
	for _, pdb := range objects.PodDisruptionBudgets.Items {
		pdb := pdb
		if pdb.Spec.Selector.MatchLabels != nil {
			tmp := true
			for _, pod := range objects.Pods.Items {
				if pdb.Namespace == pod.Namespace {
					if p.checkInUse(pod.Labels, pdb.Spec.Selector.MatchLabels) {
						tmp = true
						break
					} else {
						tmp = false
					}
				}
			}
			if !tmp && check.IsEnabled(p.Name(), &pdb.ObjectMeta) {
				d := check.Diagnostic{
					Severity: check.Warning,
					Message:  check.Message[211],
					Kind:     check.PodDisruptionBudget,
					Object:   &pdb.ObjectMeta,
					Owners:   pdb.ObjectMeta.GetOwnerReferences(),
				}
				diagnostics = append(diagnostics, d)
			}
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func (p *unusedPDBCheck) checkInUse(podLabels, pdbLabels map[string]string) bool {
	tmp := make([]int, 0)
	for k, v := range pdbLabels {
		if strings.Contains(GenerateLabelSelector(podLabels), k+"="+v) {
			tmp = append(tmp, 1)
		}
	}
	if len(tmp) == len(pdbLabels) {
		return true
	}
	return false
}

func GenerateLabelSelector(selector map[string]string) string {
	var labelSelector string
	for k, v := range selector {
		labelSelector += k + "=" + v + ","
	}
	return strings.TrimRight(labelSelector, ",")
}
