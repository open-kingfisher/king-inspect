package state

import (
	"github.com/open-kingfisher/king-inspect/check"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func initPod() *check.Objects {
	objs := &check.Objects{
		Pods: &corev1.PodList{
			Items: []corev1.Pod{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "pod_foo", Namespace: "k8s"},
				},
			},
		},
	}
	return objs
}

func GetObjectMeta() *metav1.ObjectMeta {
	objs := initPod()
	return &objs.Pods.Items[0].ObjectMeta
}

func GetOwners() []metav1.OwnerReference {
	objs := initPod()
	return objs.Pods.Items[0].ObjectMeta.GetOwnerReferences()
}
