package basic

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"kingfisher/kf/common/log"
	"kingfisher/king-inspect/check"
)

func init() {
	if err := check.Register(&metricServerCheck{}); err != nil {
		log.Errorf("Register liveness error:%s", err)
	}
}

type metricServerCheck struct{}

// Name 返回此检查的唯一名称
func (m *metricServerCheck) Name() string {
	return "metric-server"
}

// Groups 返回此检查应属于的组名列表
func (m *metricServerCheck) Groups() []string {
	return []string{"basic"}
}

// Description 返回此检查的描述信息
func (m *metricServerCheck) Description() string {
	return "检查集群Metric Server是否安装"
}

// Run 运行这个检查
func (m *metricServerCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	summary.Total = 1
	object := &metav1.ObjectMeta{
		Name: string(check.MetricServer),
	}
	if !m.clusterHasMetrics(objects) {
		d := check.Diagnostic{
			Severity: check.Suggestion,
			Message:  check.Message[120],
			Kind:     check.MetricServer,
			Object:   object,
		}
		diagnostics = append(diagnostics, d)
	}
	summary.Issue = len(diagnostics)
	summary.Suggestion = summary.Issue
	return diagnostics, summary, nil
}

func (m *metricServerCheck) clusterHasMetrics(objects *check.Objects) bool {
	supportedMetricsAPIVersions := []string{"v1beta1"}
	for _, discoveredAPIGroup := range objects.APIGroupList.Groups {
		if discoveredAPIGroup.Name != mv1beta1.GroupName {
			continue
		}
		for _, version := range discoveredAPIGroup.Versions {
			for _, supportedVersion := range supportedMetricsAPIVersions {
				if version.Version == supportedVersion {
					return true
				}
			}
		}
	}

	return false
}
