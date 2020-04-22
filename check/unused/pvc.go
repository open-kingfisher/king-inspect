package unused

import (
	"github.com/open-kingfisher/king-inspect/check"
)

func init() {
	check.Register(&unusedClaimCheck{})
}

type unusedClaimCheck struct{}

// Name 返回此检查的唯一名称
func (c *unusedClaimCheck) Name() string {
	return "unused-pvc"
}

// Groups 返回此检查应属于的组名列表
func (c *unusedClaimCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (c *unusedClaimCheck) Description() string {
	return "检查集群中未使用的PVC"
}

// Run 运行这个检查
func (c *unusedClaimCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	used := make(map[check.Identifier]struct{})
	var empty struct{}
	for _, pod := range objects.Pods.Items {
		pod := pod
		for _, volume := range pod.Spec.Volumes {
			claim := volume.VolumeSource.PersistentVolumeClaim
			if claim != nil {
				used[check.Identifier{Name: claim.ClaimName, Namespace: pod.GetNamespace()}] = empty
			}
		}
	}
	summary.Total = len(objects.PersistentVolumeClaims.Items)
	for _, claim := range objects.PersistentVolumeClaims.Items {
		claim := claim
		if _, ok := used[check.Identifier{Name: claim.GetName(), Namespace: claim.GetNamespace()}]; !ok && check.IsEnabled(c.Name(), &claim.ObjectMeta) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[201],
				Kind:     check.PersistentVolumeClaim,
				Object:   &claim.ObjectMeta,
				Owners:   claim.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
