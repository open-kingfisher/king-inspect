package unused

import (
	"github.com/open-kingfisher/king-inspect/check"
	"strings"
)

func init() {
	check.Register(&unusedPPCheck{})
}

type unusedPPCheck struct{}

// Name 返回此检查的唯一名称
func (p *unusedPPCheck) Name() string {
	return "unused-pod-preset"
}

// Groups 返回此检查应属于的组名列表
func (p *unusedPPCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (p *unusedPPCheck) Description() string {
	return "检查集群中未使用的Pod预设"
}

// Run 运行这个检查
func (p *unusedPPCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	if objects.PodPresets == nil {
		return diagnostics, summary, nil
	}
	summary.Total = len(objects.PodPresets.Items)
	for _, pp := range objects.PodPresets.Items {
		pp := pp
		if pp.Spec.Selector.MatchLabels != nil {
			tmp := true
			for _, pod := range objects.Pods.Items {
				if pp.Namespace == pod.Namespace {
					if p.checkInUse(pod.Labels, pp.Spec.Selector.MatchLabels) {
						tmp = true
						break
					} else {
						tmp = false
					}
				}
			}
			if !tmp && check.IsEnabled(p.Name(), &pp.ObjectMeta) {
				d := check.Diagnostic{
					Severity: check.Warning,
					Message:  check.Message[212],
					Kind:     check.PodPreset,
					Object:   &pp.ObjectMeta,
					Owners:   pp.ObjectMeta.GetOwnerReferences(),
				}
				diagnostics = append(diagnostics, d)
			}
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func (p *unusedPPCheck) checkInUse(podLabels, pdbLabels map[string]string) bool {
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

//func GenerateLabelSelector(selector map[string]string) string {
//	var labelSelector string
//	for k, v := range selector {
//		labelSelector += k + "=" + v + ","
//	}
//	return strings.TrimRight(labelSelector, ",")
//}
