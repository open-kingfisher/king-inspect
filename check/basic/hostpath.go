package basic

import (
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&hostPathCheck{})
}

type hostPathCheck struct{}

// Name 返回此检查的唯一名称
func (h *hostPathCheck) Name() string {
	return "hostpath-volume"
}

// Groups 返回此检查应属于的组名列表
func (h *hostPathCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (h *hostPathCheck) Description() string {
	return "检查是否有使用主机路径挂载卷的pod"
}

// Run 运行这个检查
func (h *hostPathCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		if check.IsEnabled(h.Name(), &pod.ObjectMeta) {
			for _, volume := range pod.Spec.Volumes {
				pod := pod
				if volume.VolumeSource.HostPath != nil {
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  check.Message[103],
						Kind:     check.Pod,
						Object:   &pod.ObjectMeta,
						Owners:   pod.ObjectMeta.GetOwnerReferences(),
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
