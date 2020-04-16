package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	"kingfisher/king-inspect/check"
)

func TestResourceRequestsCheckMeta(t *testing.T) {
	resourceRequirementsCheck := resourceRequirementsCheck{}
	assert.Equal(t, "resource-requirements", resourceRequirementsCheck.Name())
	assert.Equal(t, []string{"basic"}, resourceRequirementsCheck.Groups())
	assert.NotEmpty(t, resourceRequirementsCheck.Description())
}

func TestResourceRequestsCheckRegistration(t *testing.T) {
	resourceRequirementsCheck := &resourceRequirementsCheck{}
	check, err := check.Get("resource-requirements")
	assert.NoError(t, err)
	assert.Equal(t, check, resourceRequirementsCheck)
}

func TestResourceRequestsWarning(t *testing.T) {
	const message = "Set resource requests and limits for container `bar` to prevent resource contention"

	resourceRequirementsCheck := resourceRequirementsCheck{}

	tests := []struct {
		name     string
		objs     *check.Objects
		expected []check.Diagnostic
	}{
		{
			name:     "no pods",
			objs:     initPod(),
			expected: nil,
		},
		{
			name: "container with no resource requests or limits",
			objs: container("alpine"),
			expected: []check.Diagnostic{
				{
					Severity: check.Warning,
					Message:  message,
					Kind:     check.Pod,
					Object:   GetObjectMeta(),
					Owners:   GetOwners(),
				},
			},
		},
		{
			name: "init container with no resource requests or limits",
			objs: initContainer("alpine"),
			expected: []check.Diagnostic{
				{
					Severity: check.Warning,
					Message:  message,
					Kind:     check.Pod,
					Object:   GetObjectMeta(),
					Owners:   GetOwners(),
				},
			},
		},
		{
			name:     "resource requests set",
			objs:     resources(),
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := resourceRequirementsCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func resources() *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "bar",
				Image: "alpine",
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU: *resource.NewQuantity(500, "m"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU: *resource.NewQuantity(1000, "m"),
					},
				},
			}},
		InitContainers: []corev1.Container{
			{
				Name:  "bar",
				Image: "alpine",
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU: *resource.NewQuantity(500, "m"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU: *resource.NewQuantity(1000, "m"),
					},
				},
			},
		},
	}
	return objs

}
