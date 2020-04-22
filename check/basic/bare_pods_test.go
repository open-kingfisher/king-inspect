package basic

import (
	"testing"

	"github.com/open-kingfisher/king-inspect/check"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBarePodCheckMeta(t *testing.T) {
	barePodCheck := barePodCheck{}
	assert.Equal(t, "bare-pods", barePodCheck.Name())
	assert.Equal(t, []string{"basic"}, barePodCheck.Groups())
	assert.NotEmpty(t, barePodCheck.Description())
}

func TestBarePodCheckRegistration(t *testing.T) {
	barePodCheck := &barePodCheck{}
	check, err := check.Get(barePodCheck.Name())
	assert.NoError(t, err, "注册不应失败")
	assert.Equal(t, check, barePodCheck)
}

func TestBarePodError(t *testing.T) {
	tests := []struct {
		name     string
		objs     *check.Objects
		expected []check.Diagnostic
	}{
		{
			name:     "no pods",
			objs:     &check.Objects{Pods: &corev1.PodList{}},
			expected: nil,
		},
		{
			name:     "pod has owner ref",
			objs:     initRefs(initPod()),
			expected: nil,
		},
		{
			name:     "multiple pods with owner refs",
			objs:     initRefs(initMultiplePods()),
			expected: nil,
		},
		{
			name: "pod has no owner ref",
			objs: initPod(),
			expected: []check.Diagnostic{
				{
					Severity: "warning",
					Kind:     check.Pod,
					Message:  "Avoid using bare pods in clusters",
					Object:   GetObjectMeta(),
					Owners:   nil,
				},
			},
		},
		{
			name: "multiple pods with no owner ref",
			objs: initMultiplePods(),
			expected: []check.Diagnostic{
				{
					Severity: "warning",
					Kind:     check.Pod,
					Message:  "Avoid using bare pods in clusters",
					Object:   &metav1.ObjectMeta{Name: "pod_1", Namespace: "k8s"},
					Owners:   nil,
				},
				{
					Severity: "warning",
					Kind:     check.Pod,
					Message:  "Avoid using bare pods in clusters",
					Object:   &metav1.ObjectMeta{Name: "pod_2", Namespace: "k8s"},
					Owners:   nil,
				},
			},
		},
	}

	barePodCheck := &barePodCheck{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := barePodCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func initRefs(objs *check.Objects) *check.Objects {
	for index, _ := range objs.Pods.Items {
		objs.Pods.Items[index].ObjectMeta.OwnerReferences = []metav1.OwnerReference{
			{
				Name:       "Deployment",
				APIVersion: "apps/v1",
			},
		}
	}
	return objs
}
