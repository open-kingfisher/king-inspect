package state

import (
	"testing"

	"github.com/open-kingfisher/king-inspect/check"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestMeta(t *testing.T) {
	podStatusCheck := podStatusCheck{}
	assert.Equal(t, "pod-state", podStatusCheck.Name())
	assert.Equal(t, []string{"workload-health"}, podStatusCheck.Groups())
	assert.NotEmpty(t, podStatusCheck.Description())
}

func TestPodStateCheckRegistration(t *testing.T) {
	podStatusCheck := &podStatusCheck{}
	check, err := check.Get("pod-state")
	assert.NoError(t, err)
	assert.Equal(t, check, podStatusCheck)
}

func TestPodStateError(t *testing.T) {
	podStatusCheck := podStatusCheck{}
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
			name:     "pod with running status",
			objs:     status(corev1.PodRunning),
			expected: nil,
		},
		{
			name:     "pod with pending status",
			objs:     status(corev1.PodPending),
			expected: nil,
		},
		{
			name:     "pod with succeeded status",
			objs:     status(corev1.PodSucceeded),
			expected: nil,
		},
		{
			name: "pod with failed status",
			objs: status(corev1.PodFailed),
			expected: []check.Diagnostic{
				{
					Severity: check.Warning,
					Message:  "Unhealthy pod. State: `Failed`. Pod state should be `Running`, `Pending` or `Succeeded`.",
					Kind:     check.Pod,
					Object:   GetObjectMeta(),
					Owners:   GetOwners(),
				},
			},
		},
		{
			name: "pod with unknown status",
			objs: status(corev1.PodUnknown),
			expected: []check.Diagnostic{
				{
					Severity: check.Warning,
					Message:  "Unhealthy pod. State: `Unknown`. Pod state should be `Running`, `Pending` or `Succeeded`.",
					Kind:     check.Pod,
					Object:   GetObjectMeta(),
					Owners:   GetOwners(),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := podStatusCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func status(status corev1.PodPhase) *check.Objects {
	objs := initPod()
	objs.Pods.Items[0].Status = corev1.PodStatus{
		Phase: status,
	}
	return objs
}
