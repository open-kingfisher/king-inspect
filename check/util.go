package check

import (
	"golang.org/x/sync/errgroup"
	"k8s.io/api/core/v1"
	"strings"

	"sync"
)

func DeploymentReferences(objects *Objects) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, d := range objects.Deployments.Items {
		d := d
		g.Go(func() error {
			mu.Lock()
			used[Identifier{Name: d.Name, Namespace: d.Namespace}] = empty
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func StatefulSetReferences(objects *Objects) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, ss := range objects.StatefulSets.Items {
		ss := ss
		g.Go(func() error {
			mu.Lock()
			used[Identifier{Name: ss.Name, Namespace: ss.Namespace}] = empty
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func DaemonSetReferences(objects *Objects) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, ds := range objects.DaemonSets.Items {
		ds := ds
		g.Go(func() error {
			mu.Lock()
			used[Identifier{Name: ds.Name, Namespace: ds.Namespace}] = empty
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func CRBReferences(objects *Objects, kind string) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, crb := range objects.ClusterRoleBindings.Items {
		crb := crb
		g.Go(func() error {
			mu.Lock()
			for _, sub := range crb.Subjects {
				if sub.Kind == kind {
					used[Identifier{Name: sub.Name, Namespace: sub.Namespace}] = empty
				}
			}
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func RBReferences(objects *Objects, kind string) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, rb := range objects.RoleBindings.Items {
		rb := rb
		g.Go(func() error {
			mu.Lock()
			for _, sub := range rb.Subjects {
				if sub.Kind == kind {
					used[Identifier{Name: sub.Name, Namespace: sub.Namespace}] = empty
				}
			}
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func PodSAReferences(objects *Objects) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, pod := range objects.Pods.Items {
		pod := pod
		g.Go(func() error {
			mu.Lock()
			if pod.Spec.ServiceAccountName != "" {
				used[Identifier{Name: pod.Spec.ServiceAccountName, Namespace: pod.Namespace}] = empty
			}
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func SecretReferences(objects *Objects) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, s := range objects.Secrets.Items {
		s := s
		g.Go(func() error {
			mu.Lock()
			used[Identifier{Name: s.GetName(), Namespace: s.Namespace}] = empty
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func CRBRoleReferences(objects *Objects, kind string) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, crb := range objects.ClusterRoleBindings.Items {
		crb := crb
		g.Go(func() error {
			mu.Lock()
			if crb.RoleRef.Kind == kind {
				used[Identifier{Name: crb.RoleRef.Name, Namespace: ""}] = empty
			}
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func RBRoleReferences(objects *Objects, kind string) (map[Identifier]struct{}, error) {
	used := make(map[Identifier]struct{})
	var empty struct{}
	var mu sync.Mutex
	var g errgroup.Group
	for _, rb := range objects.RoleBindings.Items {
		rb := rb
		g.Go(func() error {
			mu.Lock()
			if rb.RoleRef.Kind == kind {
				used[Identifier{Name: rb.RoleRef.Name, Namespace: ""}] = empty
			}
			mu.Unlock()
			return nil
		})
	}
	return used, g.Wait()
}

func GetServerCommand(objects *Objects, server string) ([]string, v1.Pod) {
	for _, pod := range objects.KubeSystemPods.Items {
		if strings.HasPrefix(pod.Name, server) {
			if len(pod.Spec.Containers) > 0 {
				return pod.Spec.Containers[0].Command, pod
			}
		}
	}
	return []string{}, v1.Pod{}
}
