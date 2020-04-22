package basic

import (
	"fmt"
	"github.com/open-kingfisher/king-inspect/check"
	corev1 "k8s.io/api/core/v1"
	quantity "k8s.io/apimachinery/pkg/api/resource"
)

func init() {
	check.Register(&resourceRequirementsCheck{})
}

type resourceRequirementsCheck struct{}

// Name 返回此检查的唯一名称
func (r *resourceRequirementsCheck) Name() string {
	return "resource-requirements"
}

// Groups 返回此检查应属于的组名列表
func (r *resourceRequirementsCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (r *resourceRequirementsCheck) Description() string {
	return "检查Pod是否配置资源配额"
}

// Run 运行这个检查
func (r *resourceRequirementsCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		if check.IsEnabled(r.Name(), &pod.ObjectMeta) {
			d := r.checkResourceRequirements(pod.Spec.Containers, pod, "容器")
			diagnostics = append(diagnostics, d...)
			d = r.checkResourceRequirements(pod.Spec.InitContainers, pod, "初始容器")
			diagnostics = append(diagnostics, d...)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Suggestion = summary.Issue
	return diagnostics, summary, nil
}

func (r *resourceRequirementsCheck) checkResourceRequirements(containers []corev1.Container, pod corev1.Pod, kind string) []check.Diagnostic {
	var diagnostics []check.Diagnostic
	for _, container := range containers {
		container := container
		if !isLimit(container.Resources) && !isRequests(container.Resources) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  fmt.Sprintf(check.Message[110], kind, container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
		if !isLimit(container.Resources) && isRequests(container.Resources) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  fmt.Sprintf(check.Message[111], kind, container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
		// 默认设置了限制后在不设置要求的情况下，默认要求会被设置成和限制相等
		//if isLimit(container.Resources) && !isRequests(container.Resources) {
		//	d := check.Diagnostic{
		//		Severity: check.Warning,
		//		Message:  fmt.Sprintf("设置容器资源要求以防止资源争用"),
		//		Kind:     check.Pod,
		//		Object:   &pod.ObjectMeta,
		//		Owners:   pod.ObjectMeta.GetOwnerReferences(),
		//	}
		//	diagnostics = append(diagnostics, d)
		//}
		request, msg := isHighRequests(container.Resources)
		if request {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  fmt.Sprintf(msg, kind, container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics
}

func isLimit(resources corev1.ResourceRequirements) bool {
	if len(resources.Limits) > 0 {
		return true
	}
	return false
}

func isRequests(resources corev1.ResourceRequirements) bool {
	if len(resources.Requests) > 0 {
		return true
	}
	return false
}

func isHighRequests(resources corev1.ResourceRequirements) (bool, string) {
	status := 0
	if len(resources.Requests) > 0 {
		for k, v := range resources.Requests {
			// 5 core
			cpu := quantity.NewQuantity(5, quantity.DecimalSI)
			if k == corev1.ResourceCPU && v.Cmp(*cpu) == 1 {
				status += 1
			}
			// 1ki = 1024; 1Mi = 1024 * 1024
			// 5Gi
			memory := quantity.NewQuantity(5<<30, quantity.DecimalSI)
			if k == corev1.ResourceMemory && v.Cmp(*memory) == 1 {
				status += 2
			}
		}
	}
	if status == 1 {
		return true, check.Message[112]
	} else if status == 2 {
		return true, check.Message[113]
	} else if status == 3 {
		return true, check.Message[114]
	} else {
		return false, ""
	}
}
