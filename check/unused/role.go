package unused

import (
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&roleCheck{})
}

type roleCheck struct{}

// Name 返回此检查的唯一名称
func (c *roleCheck) Name() string {
	return "unused-role"
}

// Groups 返回此检查应属于的组名列表
func (c *roleCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (c *roleCheck) Description() string {
	return "检查集群中没用使用的Role"
}

// Run 运行这个检查
func (c *roleCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	used, err := check.CRBRoleReferences(objects, "Role")
	if err != nil {
		return nil, check.Summary{}, err
	}

	rb, err := check.RBRoleReferences(objects, "Role")
	if err != nil {
		return nil, check.Summary{}, err
	}

	for k, v := range rb {
		used[k] = v
	}

	var summary check.Summary
	summary.Total = len(objects.Roles.Items)
	for _, r := range objects.Roles.Items {
		if _, ok := used[check.Identifier{Name: r.GetName(), Namespace: ""}]; !ok && check.IsEnabled(c.Name(), &r.ObjectMeta) {
			r := r
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[210],
				Kind:     check.Role,
				Object:   &r.ObjectMeta,
				Owners:   r.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
