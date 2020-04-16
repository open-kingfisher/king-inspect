package security

import (
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&schedulerCheck{})
}

type schedulerCheck struct{}

// Name 返回此检查的唯一名称
func (c *schedulerCheck) Name() string {
	return "scheduler"
}

// Groups 返回此检查应属于的组名列表
func (c *schedulerCheck) Groups() []string {
	return []string{"security"}
}

// Description 返回此检查的描述信息
func (c *schedulerCheck) Description() string {
	return "检查Scheduler安全配置"
}

// Run 运行这个检查
func (c *schedulerCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	commands, pod := check.GetServerCommand(objects, "kube-scheduler-")
	if len(commands) < 1 {
		return nil, check.Summary{}, nil
	}
	result := containSchedulerCheck(commands)
	var summary check.Summary
	summary.Total = len(result)
	d := check.Diagnostic{
		Severity: check.Warning,
		Kind:     check.Scheduler,
		Object:   &pod.ObjectMeta,
	}
	for i := 0; i < len(result); i++ {
		if !result[i] {
			summary.Issue += 1
			summary.Warning += 1
			d.Message = check.Message[600+i]
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics, summary, nil
}

func containSchedulerCheck(commands []string) map[int]bool {
	result := map[int]bool{
		0: false,
		1: false,
	}
	for _, command := range commands {
		if command == "--profiling=false" {
			result[0] = true
		}
		if command == "--bind-address=127.0.0.1" {
			result[1] = true
		}
	}
	return result
}
