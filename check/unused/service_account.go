package unused

import (
	corev1 "k8s.io/api/core/v1"
	"kingfisher/king-inspect/check"
)

func init() {
	check.Register(&serviceAccountCheck{})
}

type serviceAccountCheck struct{}

// Name 返回此检查的唯一名称
func (c *serviceAccountCheck) Name() string {
	return "unused-service-account"
}

// Groups 返回此检查应属于的组名列表
func (c *serviceAccountCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (c *serviceAccountCheck) Description() string {
	return "检查集群中没用使用的ServiceAccount"
}

// Run 运行这个检查
func (c *serviceAccountCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	used, err := check.CRBReferences(objects, "ServiceAccount")
	if err != nil {
		return nil, check.Summary{}, err
	}

	rb, err := check.RBReferences(objects, "ServiceAccount")
	if err != nil {
		return nil, check.Summary{}, err
	}

	pod, err := check.PodSAReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	for k, v := range rb {
		used[k] = v
	}

	for k, v := range pod {
		used[k] = v
	}
	var summary check.Summary
	summary.Total = len(objects.ServiceAccounts.Items)
	for _, sa := range objects.ServiceAccounts.Items {
		if c.checkMounts(&sa) && check.IsEnabled(c.Name(), &sa.ObjectMeta) {
			sa := sa
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[207],
				Kind:     check.ServiceAccount,
				Object:   &sa.ObjectMeta,
				Owners:   sa.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
		//if ok, name, namespace:= c.checkSecretRefs(&sa, objects); ok {
		//	summary.Issue += 1
		//	summary.Warning += 1
		//	sa := sa
		//	d := check.Diagnostic{
		//		Severity: check.Warning,
		//		Message:  fmt.Sprintf(check.Message[208], namespace, name),
		//		Kind:     check.ServiceAccount,
		//		Object:   &sa.ObjectMeta,
		//		Owners:   sa.ObjectMeta.GetOwnerReferences(),
		//	}
		//	diagnostics = append(diagnostics, d)
		//}
		if _, ok := used[check.Identifier{Name: sa.GetName(), Namespace: sa.GetNamespace()}]; !ok && check.IsEnabled(c.Name(), &sa.ObjectMeta) {
			sa := sa
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[206],
				Kind:     check.ServiceAccount,
				Object:   &sa.ObjectMeta,
				Owners:   sa.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func (c *serviceAccountCheck) checkMounts(sa *corev1.ServiceAccount) bool {
	if sa.AutomountServiceAccountToken != nil && *sa.AutomountServiceAccountToken {
		return true
	}
	return false
}

//func (c *serviceAccountCheck) checkSecretRefs(sa *corev1.ServiceAccount, objects *check.Objects) (bool, string, string) {
//	used, _ := check.SecretReferences(objects)
//	for _, s := range sa.Secrets {
//		if s.Namespace != "" {
//			if _, ok := used[check.Identifier{Name: s.Name, Namespace: s.Namespace}]; !ok {
//				return true, s.Name, s.Namespace
//			}
//		}
//	}
//	return false, "", ""
//}
