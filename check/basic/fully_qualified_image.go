package basic

import (
	"fmt"
	"github.com/docker/distribution/reference"
	"github.com/open-kingfisher/king-inspect/check"
	"github.com/open-kingfisher/king-utils/common/log"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	if err := check.Register(&fullyQualifiedImageCheck{}); err != nil {
		log.Errorf("Register fully-qualified-image error:%s", err)
	}
}

type fullyQualifiedImageCheck struct{}

// Name 返回此检查的唯一名称
func (fq *fullyQualifiedImageCheck) Name() string {
	return "fully-qualified-image"
}

// Groups 返回此检查应属于的组名列表
func (fq *fullyQualifiedImageCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (fq *fullyQualifiedImageCheck) Description() string {
	return "检查是否使用完全合格的镜像名"
}

// Run 运行这个检查
func (fq *fullyQualifiedImageCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		if check.IsEnabled(fq.Name(), &pod.ObjectMeta) {
			d := fq.checkImage(pod.Spec.Containers, pod, "容器")
			diagnostics = append(diagnostics, d...)
			d = fq.checkImage(pod.Spec.InitContainers, pod, "初始容器")
			diagnostics = append(diagnostics, d...)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func (fq *fullyQualifiedImageCheck) checkImage(containers []corev1.Container, pod corev1.Pod, kind string) []check.Diagnostic {
	var diagnostics []check.Diagnostic
	for _, container := range containers {
		value, err := reference.ParseAnyReference(container.Image)
		if err != nil {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  fmt.Sprintf(check.Message[101], kind, container.Name, container.Image),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		} else {
			if value.String() != container.Image {
				d := check.Diagnostic{
					Severity: check.Warning,
					Message:  fmt.Sprintf(check.Message[102], kind, container.Name, container.Image),
					Kind:     check.Pod,
					Object:   &pod.ObjectMeta,
					Owners:   pod.ObjectMeta.GetOwnerReferences(),
				}
				diagnostics = append(diagnostics, d)
			}
		}
	}
	return diagnostics
}
