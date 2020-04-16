package unused

import (
	"kingfisher/king-inspect/check"
	"strings"
)

func init() {
	check.Register(&unusedServiceCheck{})
}

type unusedServiceCheck struct{}

// Name 返回此检查的唯一名称
func (p *unusedServiceCheck) Name() string {
	return "unused-service"
}

// Groups 返回此检查应属于的组名列表
func (p *unusedServiceCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (p *unusedServiceCheck) Description() string {
	return "检查集群中未使用的Service"
}

// Run 运行这个检查
func (p *unusedServiceCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary

	summary.Total = len(objects.Services.Items)
	for _, s := range objects.Services.Items {
		s := s
		if s.Spec.Selector != nil {
			tmp := true
			for _, pod := range objects.Pods.Items {
				if s.Namespace == pod.Namespace {
					if p.checkInUse(pod.Labels, s.Spec.Selector) {
						tmp = true
						break
					} else {
						tmp = false
					}
				}
			}
			if !tmp && check.IsEnabled(p.Name(), &s.ObjectMeta) {
				d := check.Diagnostic{
					Severity: check.Warning,
					Message:  check.Message[213],
					Kind:     check.Service,
					Object:   &s.ObjectMeta,
					Owners:   s.ObjectMeta.GetOwnerReferences(),
				}
				diagnostics = append(diagnostics, d)
			}
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func (p *unusedServiceCheck) checkInUse(podLabels, pdbLabels map[string]string) bool {
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
