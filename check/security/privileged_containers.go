package security

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&privilegedContainerCheck{})
}

type privilegedContainerCheck struct{}

// Name 返回此检查的唯一名称
func (pc *privilegedContainerCheck) Name() string {
	return "privileged-containers"
}

// Groups 返回此检查应属于的组名列表
func (pc *privilegedContainerCheck) Groups() []string {
	return []string{"security"}
}

// Description 返回此检查的描述信息
func (pc *privilegedContainerCheck) Description() string {
	return "检查是否有带有特权模式的容器"
}

// Run 运行这个检查
func (pc *privilegedContainerCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		diagnostics = append(diagnostics, pc.checkPrivileged(pod.Spec.Containers, pod)...)
		diagnostics = append(diagnostics, pc.checkPrivileged(pod.Spec.InitContainers, pod)...)
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

// Run 运行这个检查
func (pc *privilegedContainerCheck) checkPrivileged(containers []corev1.Container, pod corev1.Pod) []check.Diagnostic {
	var diagnostics []check.Diagnostic
	for _, container := range containers {
		if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged && check.IsEnabled(pc.Name(), &pod.ObjectMeta) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  fmt.Sprintf(check.Message[800], container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics
}
