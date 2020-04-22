package security

import (
	"github.com/open-kingfisher/king-inspect/check"
	"strings"
)

func init() {
	check.Register(&apiServerCheck{})
}

type apiServerCheck struct{}

// Name 返回此检查的唯一名称
func (c *apiServerCheck) Name() string {
	return "api-server"
}

// Groups 返回此检查应属于的组名列表
func (c *apiServerCheck) Groups() []string {
	return []string{"security"}
}

// Description 返回此检查的描述信息
func (c *apiServerCheck) Description() string {
	return "检查API Server安全配置"
}

// Run 运行这个检查
func (c *apiServerCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	commands, pod := check.GetServerCommand(objects, "kube-apiserver-")
	if len(commands) < 1 {
		return nil, check.Summary{}, nil
	}
	result := containAPIServerCheck(commands)
	var summary check.Summary
	summary.Total = len(result)
	d := check.Diagnostic{
		Severity: check.Warning,
		Kind:     check.APIServer,
		Object:   &pod.ObjectMeta,
	}
	for i := 0; i < len(result); i++ {
		if !result[i] {
			summary.Issue += 1
			summary.Warning += 1
			d.Message = check.Message[400+i]
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics, summary, nil
}

func containAPIServerCheck(commands []string) map[int]bool {
	result := map[int]bool{
		0:  false,
		1:  true,
		2:  true,
		3:  true,
		4:  false,
		5:  false,
		6:  false,
		7:  true,
		8:  false,
		9:  false,
		10: false,
		11: true,
		12: false,
		13: true,
		14: true,
		15: true,
		16: false,
		17: false,
		18: true,
		19: false,
		20: true,
		21: false,
		22: false,
		23: false,
		24: false,
		25: false,
	}
	for _, command := range commands {
		if command == "--anonymous-auth=false" {
			result[0] = true
		}
		if strings.HasPrefix(command, "--basic-auth-file") {
			result[1] = false
		}
		if strings.HasPrefix(command, "--token-auth-file") {
			result[2] = false
		}
		if command == "--kubelet-https=false" {
			result[3] = false
		}
		if strings.HasPrefix(command, "--kubelet-client-certificate") {
			result[4] = true
		}
		if strings.HasPrefix(command, "--kubelet-client-key") {
			result[5] = true
		}
		if strings.HasPrefix(command, "--kubelet-certificate-authority") {
			result[6] = true
		}
		if strings.HasPrefix(command, "--authorization-mode") {
			if strings.Contains(command, "AlwaysAllow") {
				result[7] = false
			}
			if strings.Contains(command, "Node") {
				result[8] = true
			}
			if strings.Contains(command, "RBAC") {
				result[9] = true
			}
		}
		if strings.HasPrefix(command, "--enable-admission-plugins") {
			if strings.Contains(command, "EventRateLimit") {
				result[10] = true
			}
			if strings.Contains(command, "AlwaysAdmit") {
				result[11] = false
			}
			if strings.Contains(command, "AlwaysPullImages") {
				result[12] = true
			}
			if !strings.Contains(command, "PodSecurityPolicy") && !strings.Contains(command, "SecurityContextDeny") {
				result[13] = false
			}
			if strings.Contains(command, "PodSecurityPolicy") {
				result[16] = true
			}
			if strings.Contains(command, "NodeRestriction") {
				result[17] = true
			}
		}
		if strings.HasPrefix(command, "--disable-admission-plugins") {
			if strings.Contains(command, "ServiceAccount") {
				result[14] = false
			}
			if strings.Contains(command, "NamespaceLifecycle") {
				result[15] = false
			}
		}
		if strings.HasPrefix(command, "--insecure-bind-address") {
			result[18] = false
		}
		if strings.HasPrefix(command, "--insecure-port") && command == "--insecure-port=0" {
			result[19] = true
		}
		if command == "--secure-port=0" {
			result[20] = false
		}
		if command == "--profiling=false" {
			result[21] = true
		}
		if strings.HasPrefix(command, "--audit-log-path") {
			result[22] = true
		}
		if strings.HasPrefix(command, "--audit-log-maxage") {
			result[23] = true
		}
		if strings.HasPrefix(command, "--audit-log-maxbackup") {
			result[24] = true
		}
		if strings.HasPrefix(command, "--audit-log-maxsize") {
			result[25] = true
		}
	}
	return result
}
