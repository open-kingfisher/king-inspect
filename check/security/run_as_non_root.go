package security

import (
	"fmt"
	"github.com/open-kingfisher/king-inspect/check"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	check.Register(&nonRootUserCheck{})
}

type nonRootUserCheck struct{}

// Name 返回此检查的唯一名称
func (nr *nonRootUserCheck) Name() string {
	return "non-root-user"
}

// Groups 返回此检查应属于的组名列表
func (nr *nonRootUserCheck) Groups() []string {
	return []string{"security"}
}

// Description 返回此检查的描述信息
func (nr *nonRootUserCheck) Description() string {
	return "检查是否有以根用户身份运行的pod"
}

// Run 运行这个检查
func (nr *nonRootUserCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		if check.IsEnabled(nr.Name(), &pod.ObjectMeta) {
			var containers []corev1.Container
			containers = append(containers, pod.Spec.Containers...)
			containers = append(containers, pod.Spec.InitContainers...)
			pod := pod
			podRunAsRoot := pod.Spec.SecurityContext == nil || pod.Spec.SecurityContext.RunAsNonRoot == nil || !*pod.Spec.SecurityContext.RunAsNonRoot
			for _, container := range containers {
				containerRunAsRoot := container.SecurityContext == nil || container.SecurityContext.RunAsNonRoot == nil || !*container.SecurityContext.RunAsNonRoot

				if containerRunAsRoot && podRunAsRoot {
					d := check.Diagnostic{
						Severity: check.Warning,
						Message:  fmt.Sprintf(check.Message[801], container.Name),
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
