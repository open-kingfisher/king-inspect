package unused

import (
	"sync"

	"github.com/open-kingfisher/king-inspect/check"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	check.Register(&unusedCMCheck{})
}

type unusedCMCheck struct{}

// Name 返回此检查的唯一名称
func (c *unusedCMCheck) Name() string {
	return "unused-config-map"
}

// Groups 返回此检查应属于的组名列表
func (c *unusedCMCheck) Groups() []string {
	return []string{"unused"}
}

// Description 返回此检查的描述信息
func (c *unusedCMCheck) Description() string {
	return "检查集群中没用使用的ConfigMap"
}

// Run 运行这个检查
func (c *unusedCMCheck) Run(objects *check.Objects) ([]check.Diagnostic, check.Summary, error) {
	var diagnostics []check.Diagnostic

	used, err := checkPodReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	nodeRefs, err := checkNodeReferences(objects)
	if err != nil {
		return nil, check.Summary{}, err
	}

	for k, v := range nodeRefs {
		used[k] = v
	}
	var summary check.Summary
	summary.Total = len(objects.ConfigMaps.Items)
	for _, cm := range objects.ConfigMaps.Items {
		if _, ok := used[check.Identifier{Name: cm.GetName(), Namespace: cm.GetNamespace()}]; !ok && check.IsEnabled(c.Name(), &cm.ObjectMeta) {
			cm := cm
			d := check.Diagnostic{
				Severity: check.Warning,
				Message:  check.Message[203],
				Kind:     check.ConfigMap,
				Object:   &cm.ObjectMeta,
				Owners:   cm.ObjectMeta.GetOwnerReferences(),
			}
			diagnostics = append(diagnostics, d)
		}
	}
	summary.Issue = len(diagnostics)
	summary.Warning = summary.Issue
	return diagnostics, summary, nil
}

func checkNodeReferences(objects *check.Objects) (map[check.Identifier]struct{}, error) {
	used := make(map[check.Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, node := range objects.Nodes.Items {
		node := node
		g.Go(func() error {
			source := node.Spec.ConfigSource
			if source != nil {
				mu.Lock()
				used[check.Identifier{Name: source.ConfigMap.Name, Namespace: source.ConfigMap.Namespace}] = empty
				mu.Unlock()
			}
			return nil
		})
	}
	return used, g.Wait()
}

//checkPodReferences check each pod for config map references in volumes and environment variables
func checkPodReferences(objects *check.Objects) (map[check.Identifier]struct{}, error) {
	used := make(map[check.Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, pod := range objects.Pods.Items {
		pod := pod
		namespace := pod.GetNamespace()
		g.Go(func() error {
			for _, volume := range pod.Spec.Volumes {
				cm := volume.VolumeSource.ConfigMap
				if cm != nil {
					mu.Lock()
					used[check.Identifier{Name: cm.LocalObjectReference.Name, Namespace: namespace}] = empty
					mu.Unlock()
				}
				if volume.VolumeSource.Projected != nil {
					for _, source := range volume.VolumeSource.Projected.Sources {
						cm := source.ConfigMap
						if cm != nil {
							mu.Lock()
							used[check.Identifier{Name: cm.LocalObjectReference.Name, Namespace: namespace}] = empty
							mu.Unlock()
						}
					}
				}
			}
			identifiers := checkEnvVars(pod.Spec.Containers, namespace)
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

// checkEnvVars check for config map references in container environment variables
func checkEnvVars(containers []corev1.Container, namespace string) []check.Identifier {
	var refs []check.Identifier
	for _, container := range containers {
		for _, env := range container.EnvFrom {
			if env.ConfigMapRef != nil {
				refs = append(refs, check.Identifier{Name: env.ConfigMapRef.LocalObjectReference.Name, Namespace: namespace})
			}
		}
	}
	return refs
}
