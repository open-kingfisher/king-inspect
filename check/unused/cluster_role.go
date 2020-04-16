package unused

import (
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&clusterRoleCheck{})
}

type clusterRoleCheck struct{}

// Name 返回此检查的唯一名称
func (c *clusterRoleCheck) Name() string {
	return "unused-cluster-role"
}

// Groups 返回此检查应属于的组名列表
func (c *clusterRoleCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (c *clusterRoleCheck) Description() string {
	return "检查集群中没用使用的Cluster Role"
}

// Run 运行这个检查
func (c *clusterRoleCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	used, err := check.CRBRoleReferences(objects, "ClusterRole")
	if err != nil {
		return nil, check.Summary{}, err
	}

	rb, err := check.RBRoleReferences(objects, "ClusterRole")
	if err != nil {
		return nil, check.Summary{}, err
	}

	for k, v := range rb {
		used[k] = v
	}

	var summary check.Summary
	summary.Total = len(objects.ClusterRoles.Items)
	for _, cr := range objects.ClusterRoles.Items {
		if _, ok := used[check.Identifier{Name: cr.GetName(), Namespace: ""}]; !ok && check.IsEnabled(c.Name(), &cr.ObjectMeta) {
			cr := cr
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[209],
				Kind:     check.ClusterRole,
				Object:   &cr.ObjectMeta,
				Owners:   cr.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}
