package basic

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	corev1 "k8s.io/api/core/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&latestTagCheck{})
}

type latestTagCheck struct{}

// Name 返回此检查的唯一名称
func (l *latestTagCheck) Name() string {
	return "latest-tag"
}

// Groups 返回此检查应属于的组名列表
func (l *latestTagCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (l *latestTagCheck) Description() string {
	return "检查容器镜像是否使用了latest标签"
}

// Run 运行这个检查
func (l *latestTagCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		if check.IsEnabled(l.Name(), &pod.ObjectMeta) {
			diagnostics = append(diagnostics, l.checkTags(pod.Spec.Containers, pod)...)
			diagnostics = append(diagnostics, l.checkTags(pod.Spec.InitContainers, pod)...)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue

	return diagnostics, summary, nil
}

func (l *latestTagCheck) checkTags(containers []corev1.Container, pod corev1.Pod) []check.Diagnostic {
	var diagnostics []check.Diagnostic
	for _, container := range containers {
		namedRef, err := reference.ParseNormalizedNamed(container.Image)
		if err != nil {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  fmt.Sprintf(check.Message[104], container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
			continue
		}
		tagNameOnly := reference.TagNameOnly(namedRef)
		if strings.HasSuffix(tagNameOnly.String(), ":latest") {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  fmt.Sprintf(check.Message[104], container.Name),
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics
}
