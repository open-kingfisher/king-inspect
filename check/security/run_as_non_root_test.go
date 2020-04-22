package security

import (
	"testing"

	"github.com/open-kingfisher/king-inspect/check"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestNonRootUserCheckMeta(t *testing.T) {
	nonRootUserCheck := nonRootUserCheck{}
	assert.Equal(t, "non-root-user", nonRootUserCheck.Name())
	assert.Equal(t, []string{"security"}, nonRootUserCheck.Groups())
	assert.NotEmpty(t, nonRootUserCheck.Description())
}

func TestNonRootUserCheckRegistration(t *testing.T) {
	nonRootUserCheck := &nonRootUserCheck{}
	check, err := check.Get("non-root-user")
	assert.NoError(t, err)
	assert.Equal(t, check, nonRootUserCheck)
}

func TestNonRootUserWarning(t *testing.T) {
	nonRootUserCheck := nonRootUserCheck{}
	trueVar := true
	falseVar := false

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
			name:     "pod security context and container security context unset",
			objs:     containerSecurityContextNil(),
			expected: diagnostic(),
		},
		{
			name:     "pod security context unset, container with run as non root set to true",
			objs:     containerNonRoot(nil, &trueVar),
			expected: nil,
		},
		{
			name:     "pod security context unset, container with run as non root set to false",
			objs:     containerNonRoot(nil, &falseVar),
			expected: diagnostic(),
		},
		{
			name:     "pod run as non root true, container run as non root true",
			objs:     containerNonRoot(&trueVar, &trueVar),
			expected: nil,
		},
		{
			name:     "pod run as non root true, container run as non root false",
			objs:     containerNonRoot(&trueVar, &falseVar),
			expected: nil,
		},
		{
			name:     "pod run as non root false, container run as non root true",
			objs:     containerNonRoot(&falseVar, &trueVar),
			expected: nil,
		},
		{
			name:     "pod run as non root false, container run as non root false",
			objs:     containerNonRoot(&falseVar, &falseVar),
			expected: diagnostic(),
		},
		{
			name:     "pod run as non root true, container security context unset",
			objs:     containerNonRoot(&trueVar, nil),
			expected: nil,
		},
		{
			name:     "pod run as non root false, container security context unset",
			objs:     containerNonRoot(&falseVar, nil),
			expected: diagnostic(),
		},
		// init container tests

		{
			name:     "pod security context and init container security context unset",
			objs:     initContainerSecurityContextNil(),
			expected: diagnostic(),
		},
		{
			name:     "pod security context unset, init container with run as non root set to true",
			objs:     initContainerNonRoot(nil, &trueVar),
			expected: nil,
		},
		{
			name:     "pod security context unset, init container with run as non root set to false",
			objs:     initContainerNonRoot(nil, &falseVar),
			expected: diagnostic(),
		},
		{
			name:     "pod run as non root true, init container run as non root true",
			objs:     initContainerNonRoot(&trueVar, &trueVar),
			expected: nil,
		},
		{
			name:     "pod run as non root true, init container run as non root false",
			objs:     initContainerNonRoot(&trueVar, &falseVar),
			expected: nil,
		},
		{
			name:     "pod run as non root false, init container run as non root true",
			objs:     initContainerNonRoot(&falseVar, &trueVar),
			expected: nil,
		},
		{
			name:     "pod run as non root false, init container run as non root false",
			objs:     initContainerNonRoot(&falseVar, &falseVar),
			expected: diagnostic(),
		},
		{
			name:     "pod run as non root true, init container security context unset",
			objs:     initContainerNonRoot(&trueVar, nil),
			expected: nil,
		},
		{
			name:     "pod run as non root false, init container security context unset",
			objs:     initContainerNonRoot(&falseVar, nil),
			expected: diagnostic(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := nonRootUserCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func diagnostic() []check.Diagnostic {
	pod := initPod().Pods.Items[0]
	d := []check.Diagnostic{
		{
			Severity: check.Warning,
			Message:  "Container `bar` can run as root user. Please ensure that the image is from a trusted source.",
			Kind:     check.Pod,
			Object:   &pod.ObjectMeta,
			Owners:   pod.ObjectMeta.GetOwnerReferences(),
		},
	}
	return d
}
