package state

import (
	"fmt"
	"github.com/open-kingfisher/king-inspect/check"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	check.Register(&componentStatusCheck{})
}

type componentStatusCheck struct{}

// Name 返回此检查的唯一名称
func (p *componentStatusCheck) Name() string {
	return "component-state"
}

// Groups 返回此检查应属于的组名列表
func (p *componentStatusCheck) Groups() []string {
	return []string{"state"}
}

// Description 返回此检查的描述信息
func (p *componentStatusCheck) Description() string {
	return "检查集群中不健康的Compontent"
}

// Run 运行这个检查
func (p *componentStatusCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = len(objects.Pods.Items)
	for _, component := range objects.ComponentStatuses.Items {
		component := component
		if check.IsEnabled(p.Name(), &component.ObjectMeta) {
			for _, com := range component.Conditions {
				if com.Status != corev1.ConditionTrue {
					d := check.Diagnostic{
						Severity: check.Error,
						Message:  fmt.Sprintf(check.Message[309], component.Name, com.Status, com.Message, com.Error),
						Kind:     check.Component,
						Object:   &component.ObjectMeta,
						Owners:   component.ObjectMeta.GetOwnerReferences(),
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
