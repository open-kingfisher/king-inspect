package state

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&podStatusCheck{})
}

type podStatusCheck struct{}

// Name 返回此检查的唯一名称
func (p *podStatusCheck) Name() string {
	return "pod-state"
}

// Groups 返回此检查应属于的组名列表
func (p *podStatusCheck) Groups() []string {
	return []string{"state"}
}

// Description 返回此检查的描述信息
func (p *podStatusCheck) Description() string {
	return "检查集群中不健康的Pod"
}

// Run 运行这个检查
func (p *podStatusCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	var restartCount int32
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		if check.IsEnabled(p.Name(), &pod.ObjectMeta) {
			if corev1.PodFailed == pod.Status.Phase || corev1.PodUnknown == pod.Status.Phase || corev1.PodPending == pod.Status.Phase {
				d := check.Diagnostic{
					Severity: check.Warning,
					Message:  fmt.Sprintf(check.Message[307], pod.Status.Phase),
					Kind:     check.Pod,
					Object:   &pod.ObjectMeta,
					Owners:   pod.ObjectMeta.GetOwnerReferences(),
				}
				diagnostics = append(diagnostics, d)
			}
			for _, container := range pod.Status.ContainerStatuses {
				restartCount = 10
				if container.RestartCount >= restartCount {
					d := check.Diagnostic{
						Severity: check.Suggestion,
						Message:  fmt.Sprintf(check.Message[308], container.Name, container.RestartCount, restartCount),
						Kind:     check.Pod,
						Object:   &pod.ObjectMeta,
						Owners:   pod.ObjectMeta.GetOwnerReferences(),
					}
					diagnostics = append(diagnostics, d)
				}
			}
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
