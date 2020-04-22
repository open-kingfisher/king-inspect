package security

import (
	"github.com/open-kingfisher/king-inspect/check"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivilegedContainersCheckMeta(t *testing.T) {
	privilegedContainerCheck := privilegedContainerCheck{}
	assert.Equal(t, "privileged-containers", privilegedContainerCheck.Name())
	assert.Equal(t, []string{"security"}, privilegedContainerCheck.Groups())
	assert.NotEmpty(t, privilegedContainerCheck.Description())
}

func TestPrivilegedContainersCheckRegistration(t *testing.T) {
	privilegedContainerCheck := &privilegedContainerCheck{}
	check, err := check.Get("privileged-containers")
	assert.NoError(t, err)
	assert.Equal(t, check, privilegedContainerCheck)
}

func TestPrivilegedContainerWarning(t *testing.T) {
	privilegedContainerCheck := privilegedContainerCheck{}

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
			name:     "pod with container in privileged mode",
			objs:     containerPrivileged(true),
			expected: warnings(containerPrivileged(true), privilegedContainerCheck.Name()),
		},
		{
			name:     "pod with container.SecurityContext = nil",
			objs:     containerSecurityContextNil(),
			expected: nil,
		},
		{
			name:     "pod with container.SecurityContext.Privileged = nil",
			objs:     containerPrivilegedNil(),
			expected: nil,
		},
		{
			name:     "pod with container in regular mode",
			objs:     containerPrivileged(false),
			expected: nil,
		},
		{
			name:     "pod with init container in privileged mode",
			objs:     initContainerPrivileged(true),
			expected: warnings(initContainerPrivileged(true), privilegedContainerCheck.Name()),
		},
		{
			name:     "pod with initContainer.SecurityContext = nil",
			objs:     initContainerSecurityContextNil(),
			expected: nil,
		},
		{
			name:     "pod with initContainer.SecurityContext.Privileged = nil",
			objs:     initContainerPrivilegedNil(),
			expected: nil,
		},
		{
			name:     "pod with init container in regular mode",
			objs:     initContainerPrivileged(false),
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := privilegedContainerCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func warnings(objs *check.Objects, name string) []check.Diagnostic {
	pod := objs.Pods.Items[0]
	d := []check.Diagnostic{
		{
			Severity: check.Warning,
			Message:  "Privileged container 'bar' found. Please ensure that the image is from a trusted source.",
			Kind:     check.Pod,
			Object:   &pod.ObjectMeta,
			Owners:   pod.ObjectMeta.GetOwnerReferences(),
		},
	}
	return d
}
