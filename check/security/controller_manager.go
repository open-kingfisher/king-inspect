package security

import (
	"github.com/open-kingfisher/king-inspect/check"
	"strings"
)

func init() {
	check.Register(&controllerManagerCheck{})
}

type controllerManagerCheck struct{}

// Name 返回此检查的唯一名称
func (c *controllerManagerCheck) Name() string {
	return "controller-manager"
}

// Groups 返回此检查应属于的组名列表
func (c *controllerManagerCheck) Groups() []string {
	return []string{"security"}
}

// Description 返回此检查的描述信息
func (c *controllerManagerCheck) Description() string {
	return "检查Controller Manager安全配置"
}

// Run 运行这个检查
func (c *controllerManagerCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	commands, pod := check.GetServerCommand(objects, "kube-controller-manager-")
	if len(commands) < 1 {
		return nil, check.Summary{}, nil
	}
	result := containCMCheck(commands)
	var summary check.Summary
	summary.Total = len(result)
	d := check.Diagnostic{
		Severity: check.Warning,
		Kind:     check.ControllerManager,
		Object:   &pod.ObjectMeta,
	}
	for i := 0; i < len(result); i++ {
		if !result[i] {
			summary.Issue += 1
			summary.Warning += 1
			d.Message = check.Message[500+i]
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics, summary, nil
}

func containCMCheck(commands []string) map[int]bool {
	result := map[int]bool{
		0: false,
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
		6: false,
	}
	for _, command := range commands {
		if strings.HasPrefix(command, "--terminated-pod-gc-threshold") {
			result[0] = true
		}
		if command == "--profiling=false" {
			result[1] = true
		}
		if command == "--use-service-account-credentials=true" {
			result[2] = true
		}
		if strings.HasPrefix(command, "--service-account-private-key-file") {
			result[3] = true
		}
		if strings.HasPrefix(command, "--root-ca-file") {
			result[4] = true
		}
		if strings.HasPrefix(command, "RotateKubeletServerCertificate=true") {
			result[5] = true
		}
		if command == "--bind-address=127.0.0.1" {
			result[6] = true
		}
	}
	return result
}
