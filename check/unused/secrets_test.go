package unused

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kingfisher/king-inspect/check"
)

func TestUnusedSecretCheckMeta(t *testing.T) {
	unusedSecretCheck := unusedSecretCheck{}
	assert.Equal(t, "unused-secret", unusedSecretCheck.Name())
	assert.Equal(t, []string{"basic"}, unusedSecretCheck.Groups())
	assert.NotEmpty(t, unusedSecretCheck.Description())
}

func TestUnusedSecretCheckRegistration(t *testing.T) {
	unusedSecretCheck := &unusedSecretCheck{}
	check, err := check.Get("unused-secret")
	assert.NoError(t, err)
	assert.Equal(t, check, unusedSecretCheck)
}

func TestUnusedSecretWarning(t *testing.T) {
	unusedSecretCheck := unusedSecretCheck{}

	tests := []struct {
		name     string
		objs     *check.Objects
		expected []check.Diagnostic
	}{
		{
			name:     "no secrets",
			objs:     &check.Objects{Pods: &corev1.PodList{}, Secrets: &corev1.SecretList{}},
			expected: nil,
		},
		{
			name:     "secret volume",
			objs:     secretVolume(),
			expected: nil,
		},
		{
			name:     "environment variable references secret",
			objs:     secretEnvSource(),
			expected: nil,
		},
		{
			name:     "pod with image pull secrets",
			objs:     imagePullSecrets(),
			expected: nil,
		},
		{
			name:     "projected volume references secret",
			objs:     secretProjection(),
			expected: nil,
		},
		{
			name: "unused secret",
			objs: initSecret(),
			expected: []check.Diagnostic{
				{
					Severity: check.Warning,
					Message:  "Unused secret",
					Kind:     check.Secret,
					Object:   &metav1.ObjectMeta{Name: "secret_foo", Namespace: "k8s"},
					Owners:   GetOwners(),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := unusedSecretCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func initSecret() *check.Objects {
	objs := &check.Objects{
		Pods: &corev1.PodList{
			Items: []corev1.Pod{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "pod_foo", Namespace: "k8s"},
				},
			},
		},
		Secrets: &corev1.SecretList{
			Items: []corev1.Secret{
				{
					TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{Name: "secret_foo", Namespace: "k8s"},
				},
			},
		},
	}
	return objs
}

func secretVolume() *check.Objects {
	objs := initSecret()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: "bar",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "secret_foo",
					},
				},
			}},
	}
	return objs
}

func secretProjection() *check.Objects {
	objs := initSecret()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: "bar",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{
						Sources: []corev1.VolumeProjection{
							{
								Secret: &corev1.SecretProjection{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "secret_foo",
									},
								},
							},
						},
					},
				},
			}},
	}
	return objs
}

func secretEnvSource() *check.Objects {
	objs := initSecret()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "test-container",
				Image: "docker.io/nginx",
				EnvFrom: []corev1.EnvFromSource{
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: "secret_foo"},
						},
					},
				},
			}},
	}
	return objs
}

func imagePullSecrets() *check.Objects {
	objs := initSecret()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		ImagePullSecrets: []corev1.LocalObjectReference{
			{
				Name: "secret_foo",
			},
		},
	}
	return objs
}
