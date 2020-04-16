package basic

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"kingfisher/kf/common/log"
	"kingfisher/king-inspect/check"
)

func init() {
	if err := check.Register(&livenessCheck{}); err != nil {
		log.Errorf("Register liveness error:%s", err)
	}
}

type livenessCheck struct{}

// Name 返回此检查的唯一名称
func (l *livenessCheck) Name() string {
	return "liveness"
}

// Groups 返回此检查应属于的组名列表
func (l *livenessCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (l *livenessCheck) Description() string {
	return "检查存活探针"
}

// Run 运行这个检查
func (l *livenessCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		d := l.checkLivenessProbe(pod.Spec.Containers, pod)
		diagnostics = append(diagnostics, d...)
		//d = l.checkLivenessProbe(pod.Spec.InitContainers, pod)
		//diagnostics = append(diagnostics, d...)
	}
	summary.Issue = len(diagnostics)
	summary.Suggestion = summary.Issue
	return diagnostics, summary, nil
}

func (l *livenessCheck) checkLivenessProbe(containers []corev1.Container, pod corev1.Pod) []check.Diagnostic {
	var diagnostics []check.Diagnostic
	for _, container := range containers {
		if container.LivenessProbe == nil && check.IsEnabled(l.Name(), &pod.ObjectMeta) {
			d := check.Diagnostic{
				Severity: check.Suggestion,
				Message:  fmt.Sprintf(check.Message[106], container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics
}
