package security

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

func containerPrivileged(privileged bool) *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:            "bar",
				SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
			}},
	}
	return objs
}

func containerNonRoot(pod, container *bool) *check.Objects {
	objs := initPod()
	podSecurityContext := &corev1.PodSecurityContext{}
	if pod != nil {
		podSecurityContext = &corev1.PodSecurityContext{RunAsNonRoot: pod}
	}
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		SecurityContext: podSecurityContext,
		Containers: []corev1.Container{
			{
				Name:            "bar",
				SecurityContext: &corev1.SecurityContext{RunAsNonRoot: container},
			}},
	}
	return objs
}

func initContainerNonRoot(pod, container *bool) *check.Objects {
	objs := initPod()
	podSecurityContext := &corev1.PodSecurityContext{}
	if pod != nil {
		podSecurityContext = &corev1.PodSecurityContext{RunAsNonRoot: pod}
	}
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		SecurityContext: podSecurityContext,
		InitContainers: []corev1.Container{
			{
				Name:            "bar",
				SecurityContext: &corev1.SecurityContext{RunAsNonRoot: container},
			}},
	}
	return objs
}

func containerSecurityContextNil() *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name: "bar",
			}},
	}
	return objs
}

func containerPrivilegedNil() *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:            "bar",
				SecurityContext: &corev1.SecurityContext{},
			}},
	}
	return objs
}

func initContainerPrivileged(privileged bool) *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		InitContainers: []corev1.Container{
			{
				Name:            "bar",
				SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
			}},
	}
	return objs
}

func initContainerSecurityContextNil() *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		InitContainers: []corev1.Container{
			{
				Name: "bar",
			}},
	}
	return objs
}

func initContainerPrivilegedNil() *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		InitContainers: []corev1.Container{
			{
				Name:            "bar",
				SecurityContext: &corev1.SecurityContext{},
			}},
	}
	return objs
}
