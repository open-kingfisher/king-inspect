package unused

import (
	"sync"

	"github.com/open-kingfisher/king-inspect/check"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	check.Register(&unusedSecretCheck{})
}

type unusedSecretCheck struct{}

type identifier struct {
	Name      string
	Namespace string
}

// Name 返回此检查的唯一名称
func (s *unusedSecretCheck) Name() string {
	return "unused-secret"
}

// Groups 返回此检查应属于的组名列表
func (s *unusedSecretCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (s *unusedSecretCheck) Description() string {
	return "检查集群中没用使用的Secret. 忽略 service account tokens"
}

// Run 运行这个检查
func (s *unusedSecretCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic
	var summary check.Summary
	used, err := checkReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}
	summary.Total = len(objects.Secrets.Items)
	for _, secret := range filter(objects.Secrets.Items) {
		secret := secret
		if _, ok := used[check.Identifier{Name: secret.GetName(), Namespace: secret.GetNamespace()}]; !ok && check.IsEnabled(s.Name(), &secret.ObjectMeta) {
			secret := secret
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[202],
				Kind:     check.Secret,
				Object:   &secret.ObjectMeta,
				Owners:   secret.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

//checkReferences check each pod for config map references in volumes and environment variables
func checkReferences(objects *check.Objects) (map[check.Identifier]struct{}, error) {
	used := make(map[check.Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, pod := range objects.Pods.Items {
		pod := pod
		namespace := pod.GetNamespace()
		g.Go(func() error {
			for _, volume := range pod.Spec.Volumes {
				s := volume.VolumeSource.Secret
				if s != nil {
					mu.Lock()
					used[check.Identifier{Name: s.SecretName, Namespace: namespace}] = empty
					mu.Unlock()
				}
				if volume.VolumeSource.Projected != nil {
					for _, source := range volume.VolumeSource.Projected.Sources {
						s := source.Secret
						if s != nil {
							mu.Lock()
							used[check.Identifier{Name: s.LocalObjectReference.Name, Namespace: namespace}] = empty
							mu.Unlock()
						}
					}
				}
			}
			for _, imageSecret := range pod.Spec.ImagePullSecrets {
				mu.Lock()
				used[check.Identifier{Name: imageSecret.Name, Namespace: namespace}] = empty
				mu.Unlock()
			}
			identifiers := envVarsSecretRefs(pod.Spec.Containers, namespace)
			identifiers = append(identifiers, checkEnvVars(pod.Spec.InitContainers, namespace)...)
			mu.Lock()
			for _, i := range identifiers {
				used[i] = empty
			}
			mu.Unlock()

			return nil
		})
	}

	return used, g.Wait()
}

// envVarsSecretRefs check for config map references in container environment variables
func envVarsSecretRefs(containers []corev1.Container, namespace string) []check.Identifier {
	var refs []check.Identifier
	for _, container := range containers {
		for _, env := range container.EnvFrom {
			if env.SecretRef != nil {
				refs = append(refs, check.Identifier{Name: env.SecretRef.LocalObjectReference.Name, Namespace: namespace})
			}
		}
		for _, env := range container.Env {
			if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
				refs = append(refs, check.Identifier{Name: env.ValueFrom.SecretKeyRef.LocalObjectReference.Name, Namespace: namespace})
			}
		}
	}
	return refs
}

// filter returns Secrets that are not of type `checkrnetes.io/service-account-token`
func filter(secrets []corev1.Secret) []corev1.Secret {
	var filtered []corev1.Secret
	for _, secret := range secrets {
		if secret.Type != corev1.SecretTypeServiceAccountToken {
			filtered = append(filtered, secret)
		}
	}
	return filtered
}
