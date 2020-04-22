package basic

import (
	"github.com/open-kingfisher/king-inspect/check"
	"github.com/open-kingfisher/king-utils/common/log"
)

func init() {
	if err := check.Register(&barePodCheck{}); err != nil {
		log.Errorf("Register bare-pod error:%s", err)
	}
}

type barePodCheck struct{}

// Name 返回此检查的唯一名称
func (b *barePodCheck) Name() string {
	return "bare-pod"
}

// Groups 返回此检查应属于的组名列表
func (b *barePodCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (b *barePodCheck) Description() string {
	return "检查集群中是否有裸Pod"
}

// Run 运行这个检查
func (b *barePodCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, pod := range objects.Pods.Items {
		pod := pod
		if len(pod.ObjectMeta.OwnerReferences) == 0 && check.IsEnabled(b.Name(), &pod.ObjectMeta) {
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[100],
				Kind:     check.Pod,
				Object:   &pod.ObjectMeta,
				Owners:   pod.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
