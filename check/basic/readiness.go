package basic

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"kingfisher/kf/common/log"
	"kingfisher/king-inspect/check"
)

func init() {
	if err := check.Register(&readinessCheck{}); err != nil {
		log.Errorf("Register readiness error:%s", err)
	}
}

type readinessCheck struct{}

// Name 返回此检查的唯一名称
func (r *readinessCheck) Name() string {
	return "readiness"
}

// Groups 返回此检查应属于的组名列表
func (r *readinessCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (r *readinessCheck) Description() string {
	return "检查就绪探针"
}

// Run 运行这个检查
func (r *readinessCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		if check.IsEnabled(r.Name(), &pod.ObjectMeta) {
			d := r.checkReadinessProbe(pod.Spec.Containers, pod)
			diagnostics = append(diagnostics, d...)
			//d = r.checkReadinessProbe(pod.Spec.InitContainers, pod)
			//diagnostics = append(diagnostics, d...)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Suggestion = summary.Issue
	return diagnostics, summary, nil
}

func (r *readinessCheck) checkReadinessProbe(containers []corev1.Container, pod corev1.Pod) []check.Diagnostic {
	var diagnostics []check.Diagnostic
	for _, container := range containers {
		if container.ReadinessProbe == nil {
			d := check.Diagnostic{
				Severity: check.Suggestion,
				Message:  fmt.Sprintf(check.Message[109], container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics
}
