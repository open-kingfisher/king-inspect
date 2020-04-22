package security

import (
	"github.com/open-kingfisher/king-inspect/check"
	"strings"
)

func init() {
	check.Register(&etcdCheck{})
}

type etcdCheck struct{}

// Name 返回此检查的唯一名称
func (c *etcdCheck) Name() string {
	return "etcd"
}

// Groups 返回此检查应属于的组名列表
func (c *etcdCheck) Groups() []string {
	return []string{"security"}
}

// Description 返回此检查的描述信息
func (c *etcdCheck) Description() string {
	return "检查ETCD安全配置"
}

// Run 运行这个检查
func (c *etcdCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	commands, pod := check.GetServerCommand(objects, "etcd-")
	if len(commands) < 1 {
		return nil, check.Summary{}, nil
	}
	result := containEtcdCheck(commands)
	var summary check.Summary
	summary.Total = len(result)
	d := check.Diagnostic{
		Severity: check.Warning,
		Kind:     check.ETCD,
		Object:   &pod.ObjectMeta,
	}
	for i := 0; i < len(result); i++ {
		if !result[i] {
			summary.Issue += 1
			summary.Warning += 1
			d.Message = check.Message[700+i]
			diagnostics = append(diagnostics, d)
		}
	}
	return diagnostics, summary, nil
}

func containEtcdCheck(commands []string) map[int]bool {
	result := map[int]bool{
		0: false,
		1: false,
		2: false,
		3: true,
		4: false,
		5: false,
		6: false,
		7: true,
		8: false,
	}
	for _, command := range commands {
		if strings.HasPrefix(command, "--cert-file") {
			result[0] = true
		}
		if strings.HasPrefix(command, "--key-file") {
			result[1] = true
		}
		if command == "--client-cert-auth=true" {
			result[2] = true
		}
		if command == "--auto-tls=true" {
			result[3] = false
		}
		if strings.HasPrefix(command, "--peer-cert-file") {
			result[4] = true
		}
		if strings.HasPrefix(command, "--peer-key-file") {
			result[5] = true
		}
		if command == "--peer-client-cert-auth=true" {
			result[6] = true
		}
		if command == "--peer-auto-tls=true" {
			result[7] = false
		}
		if strings.HasPrefix(command, "--trusted-ca-file") {
			result[8] = true
		}
	}
	return result
}
