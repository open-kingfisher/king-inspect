package unused

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kingfisher/king-inspect/check"
)

const cmNamespace = "k8s"

func TestUnusedConfigMapCheckMeta(t *testing.T) {
	unusedCMCheck := unusedCMCheck{}
	assert.Equal(t, "unused-config-map", unusedCMCheck.Name())
	assert.Equal(t, []string{"basic"}, unusedCMCheck.Groups())
	assert.NotEmpty(t, unusedCMCheck.Description())
}

func TestUnusedConfigMapCheckRegistration(t *testing.T) {
	unusedCMCheck := &unusedCMCheck{}
	check, err := check.Get("unused-config-map")
	assert.NoError(t, err)
	assert.Equal(t, check, unusedCMCheck)
}

func TestUnusedConfigMapWarning(t *testing.T) {
	unusedCMCheck := unusedCMCheck{}

	tests := []struct {
		name     string
		objs     *check.Objects
		expected []check.Diagnostic
	}{
		{
			name:     "no config maps",
			objs:     &check.Objects{Nodes: &corev1.NodeList{}, Pods: &corev1.PodList{}, ConfigMaps: &corev1.ConfigMapList{}},
			expected: nil,
		},
		{
			name:     "volume mounted config map",
			objs:     configMapVolume(),
			expected: nil,
		},
		{
			name:     "environment variable references config map",
			objs:     configMapEnvSource(),
			expected: nil,
		},
		{
			name:     "projected volume references config map",
			objs:     projectedVolume(),
			expected: nil,
		},
		{
			name:     "node config source references config map",
			objs:     nodeConfigSource(),
			expected: nil,
		},
		{
			name: "unused config map",
			objs: initConfigMap(),
			expected: []check.Diagnostic{
				{
					Severity: check.Warning,
					Message:  "Unused config map",
					Kind:     check.ConfigMap,
					Object:   &metav1.ObjectMeta{Name: "cm_foo", Namespace: cmNamespace},
					Owners:   GetOwners(),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := unusedCMCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func initConfigMap() *check.Objects {
	objs := &check.Objects{
		Nodes: &corev1.NodeList{
			Items: []corev1.Node{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Node", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "node_foo"},
				},
			},
		},
		Pods: &corev1.PodList{
			Items: []corev1.Pod{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "pod_foo", Namespace: cmNamespace},
				},
			},
		},
		ConfigMaps: &corev1.ConfigMapList{
			Items: []corev1.ConfigMap{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "cm_foo", Namespace: cmNamespace},
				},
			},
		},
	}
	return objs
}

func nodeConfigSource() *check.Objects {
	objs := initConfigMap()
	objs.Nodes.Items[0].Spec = corev1.NodeSpec{
		ConfigSource: &corev1.NodeConfigSource{
			ConfigMap: &corev1.ConfigMapNodeConfigSource{
				Name:      "cm_foo",
				Namespace: cmNamespace,
			},
		},
	}
	return objs
}

func configMapVolume() *check.Objects {
	objs := initConfigMap()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: "bar",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: "cm_foo"},
					},
				},
			}},
	}
	return objs
}

func projectedVolume() *check.Objects {
	objs := initConfigMap()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: "bar",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{
						Sources: []corev1.VolumeProjection{
							{
								ConfigMap: &corev1.ConfigMapProjection{
									LocalObjectReference: corev1.LocalObjectReference{Name: "cm_foo"},
								},
							},
						},
					},
				},
			}},
	}
	return objs
}

func configMapEnvSource() *check.Objects {
	objs := initConfigMap()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "test-container",
				Image: "docker.io/nginx",
				EnvFrom: []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: "cm_foo"},
						},
					},
				},
			}},
	}
	return objs
}
