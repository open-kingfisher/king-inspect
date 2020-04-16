package unused

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kingfisher/king-inspect/check"
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

func initMultiplePods() *check.Objects {
	objs := &check.Objects{
		Pods: &corev1.PodList{
			Items: []corev1.Pod{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "pod_1", Namespace: "k8s"},
				},
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "pod_2", Namespace: "k8s"},
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

func container(image string) *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "bar",
				Image: image,
			}},
	}
	return objs
}

func initContainer(image string) *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		InitContainers: []corev1.Container{
			{
				Name:  "bar",
				Image: image,
			}},
	}
	return objs
}

func issues(severity check.Severity, message string, kind check.Kind, checks string) []check.Diagnostic {
	d := []check.Diagnostic{
		{
			Severity: severity,
			Message:  message,
			Kind:     kind,
			Object:   GetObjectMeta(),
			Owners:   GetOwners(),
		},
	}
	return d
}
