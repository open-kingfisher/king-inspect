package check

import (
	"context"
	"golang.org/x/sync/errgroup"
	ar "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/api/apps/v1"
	hpav1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	pv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/api/rbac/v1beta1"
	v1alpha1 "k8s.io/api/settings/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Identifier 用于标识特定的Namespace作用域对象
type Identifier struct {
	Name      string
	Namespace string
}

// Objects 封装了Kubernetes集群中的所有对象
type Objects struct {
	Nodes                           *corev1.NodeList
	PersistentVolumes               *corev1.PersistentVolumeList
	ComponentStatuses               *corev1.ComponentStatusList
	SystemNamespace                 *corev1.Namespace
	Pods                            *corev1.PodList
	KubeSystemPods                  *corev1.PodList
	PodTemplates                    *corev1.PodTemplateList
	PersistentVolumeClaims          *corev1.PersistentVolumeClaimList
	ConfigMaps                      *corev1.ConfigMapList
	Services                        *corev1.ServiceList
	Secrets                         *corev1.SecretList
	ServiceAccounts                 *corev1.ServiceAccountList
	ResourceQuotas                  *corev1.ResourceQuotaList
	LimitRanges                     *corev1.LimitRangeList
	MutatingWebhookConfigurations   *ar.MutatingWebhookConfigurationList
	ValidatingWebhookConfigurations *ar.ValidatingWebhookConfigurationList
	Namespaces                      *corev1.NamespaceList
	Deployments                     *v1.DeploymentList
	ReplicaSets                     *v1.ReplicaSetList
	StatefulSets                    *v1.StatefulSetList
	DaemonSets                      *v1.DaemonSetList
	HPA                             *hpav1.HorizontalPodAutoscalerList
	ClusterRoleBindings             *v1beta1.ClusterRoleBindingList
	RoleBindings                    *v1beta1.RoleBindingList
	ClusterRoles                    *v1beta1.ClusterRoleList
	Roles                           *v1beta1.RoleList
	PodDisruptionBudgets            *pv1beta1.PodDisruptionBudgetList
	APIGroupList                    *metav1.APIGroupList
	PodPresets                      *v1alpha1.PodPresetList
}

// FetchObjects 从Kubernetes集群返回对象
// ctx is currently unused during API calls. More info: https://github.com/kubernetes/community/pull/1166
func FetchObjects(clientSet *kubernetes.Clientset, namespace []string, ctx context.Context) (*Objects, error) {
	client := clientSet.CoreV1()
	admissionControllerClient := clientSet.AdmissionregistrationV1beta1()
	opts := metav1.ListOptions{}
	objects := &Objects{}

	var g errgroup.Group

	g.Go(func() (err error) {
		objects.Nodes, err = client.Nodes().List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.PersistentVolumes, err = client.PersistentVolumes().List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.ComponentStatuses, err = client.ComponentStatuses().List(opts)
		return
	})
	g.Go(func() (err error) {
		if all, err := client.Pods(corev1.NamespaceAll).List(opts); err == nil {
			objects.Pods = &corev1.PodList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.Pods.Items = append(objects.Pods.Items, d)
					}
				}
			}
		}
		// 这种方式太慢要发起多次请求
		//for _, n := range namespace {
		//	all, _ := client.Pods(n).List(opts)
		//	objects.Pods.Items = append(objects.Pods.Items, all.Items...)
		//}
		return
	})
	g.Go(func() (err error) {
		if all, err := client.PodTemplates(corev1.NamespaceAll).List(opts); err == nil {
			objects.PodTemplates = &corev1.PodTemplateList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.PodTemplates.Items = append(objects.PodTemplates.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		objects.KubeSystemPods, err = client.Pods(metav1.NamespaceSystem).List(opts)
		return
	})
	g.Go(func() (err error) {
		if all, err := client.PersistentVolumeClaims(corev1.NamespaceAll).List(opts); err == nil {
			objects.PersistentVolumeClaims = &corev1.PersistentVolumeClaimList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.PersistentVolumeClaims.Items = append(objects.PersistentVolumeClaims.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := client.ConfigMaps(corev1.NamespaceAll).List(opts); err == nil {
			objects.ConfigMaps = &corev1.ConfigMapList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.ConfigMaps.Items = append(objects.ConfigMaps.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := client.Secrets(corev1.NamespaceAll).List(opts); err == nil {
			objects.Secrets = &corev1.SecretList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.Secrets.Items = append(objects.Secrets.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := client.Services(corev1.NamespaceAll).List(opts); err == nil {
			objects.Services = &corev1.ServiceList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.Services.Items = append(objects.Services.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := client.ServiceAccounts(corev1.NamespaceAll).List(opts); err == nil {
			objects.ServiceAccounts = &corev1.ServiceAccountList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.ServiceAccounts.Items = append(objects.ServiceAccounts.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := client.ResourceQuotas(corev1.NamespaceAll).List(opts); err == nil {
			objects.ResourceQuotas = &corev1.ResourceQuotaList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.ResourceQuotas.Items = append(objects.ResourceQuotas.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := client.LimitRanges(corev1.NamespaceAll).List(opts); err == nil {
			objects.LimitRanges = &corev1.LimitRangeList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.LimitRanges.Items = append(objects.LimitRanges.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		objects.SystemNamespace, err = client.Namespaces().Get(metav1.NamespaceSystem, metav1.GetOptions{})
		return
	})
	g.Go(func() (err error) {
		objects.MutatingWebhookConfigurations, err = admissionControllerClient.MutatingWebhookConfigurations().List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.ValidatingWebhookConfigurations, err = admissionControllerClient.ValidatingWebhookConfigurations().List(opts)
		return
	})
	g.Go(func() (err error) {
		objects.Namespaces, err = client.Namespaces().List(opts)
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.AppsV1().Deployments(corev1.NamespaceAll).List(opts); err == nil {
			objects.Deployments = &v1.DeploymentList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.Deployments.Items = append(objects.Deployments.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.AppsV1().StatefulSets(corev1.NamespaceAll).List(opts); err == nil {
			objects.StatefulSets = &v1.StatefulSetList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.StatefulSets.Items = append(objects.StatefulSets.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.AppsV1().ReplicaSets(corev1.NamespaceAll).List(opts); err == nil {
			objects.ReplicaSets = &v1.ReplicaSetList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.ReplicaSets.Items = append(objects.ReplicaSets.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.AppsV1().DaemonSets(corev1.NamespaceAll).List(opts); err == nil {
			objects.DaemonSets = &v1.DaemonSetList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.DaemonSets.Items = append(objects.DaemonSets.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.AutoscalingV1().HorizontalPodAutoscalers(corev1.NamespaceAll).List(opts); err == nil {
			objects.HPA = &hpav1.HorizontalPodAutoscalerList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.HPA.Items = append(objects.HPA.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		objects.ClusterRoleBindings, err = clientSet.RbacV1beta1().ClusterRoleBindings().List(opts)
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.RbacV1beta1().RoleBindings(corev1.NamespaceAll).List(opts); err == nil {
			objects.RoleBindings = &v1beta1.RoleBindingList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.RoleBindings.Items = append(objects.RoleBindings.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		objects.ClusterRoles, err = clientSet.RbacV1beta1().ClusterRoles().List(opts)
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.RbacV1beta1().Roles(corev1.NamespaceAll).List(opts); err == nil {
			objects.Roles = &v1beta1.RoleList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.Roles.Items = append(objects.Roles.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.PolicyV1beta1().PodDisruptionBudgets(corev1.NamespaceAll).List(opts); err == nil {
			objects.PodDisruptionBudgets = &pv1beta1.PodDisruptionBudgetList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.PodDisruptionBudgets.Items = append(objects.PodDisruptionBudgets.Items, d)
					}
				}
			}
		}
		return
	})
	g.Go(func() (err error) {
		objects.APIGroupList, err = clientSet.Discovery().ServerGroups()
		return
	})
	g.Go(func() (err error) {
		if all, err := clientSet.SettingsV1alpha1().PodPresets(corev1.NamespaceAll).List(opts); err == nil {
			objects.PodPresets = &v1alpha1.PodPresetList{}
			for _, d := range all.Items {
				for _, n := range namespace {
					if d.Namespace == n {
						objects.PodPresets.Items = append(objects.PodPresets.Items, d)
					}
				}
			}
		}
		return
	})
	err := g.Wait()
	if err != nil {
		return nil, err
	}

	return objects, nil
}
