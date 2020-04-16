package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kingfisher/king-inspect/check"
)

func TestFullyQualifiedImageCheckMeta(t *testing.T) {
	fullyQualifiedImageCheck := fullyQualifiedImageCheck{}
	assert.Equal(t, "fully-qualified-image", fullyQualifiedImageCheck.Name())
	assert.Equal(t, []string{"basic"}, fullyQualifiedImageCheck.Groups())
	assert.NotEmpty(t, fullyQualifiedImageCheck.Description())
}

func TestFullyQualifiedImageCheckRegistration(t *testing.T) {
	fullyQualifiedImageCheck := &fullyQualifiedImageCheck{}
	check, err := check.Get("fully-qualified-image")
	assert.NoError(t, err)
	assert.Equal(t, check, fullyQualifiedImageCheck)
}

func TestFullyQualifiedImageWarning(t *testing.T) {
	const message = "Use fully qualified image for container 'bar'"
	const severity = check.Warning
	const name = "fully-qualified-image"

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
			name:     "pod with container image - k8s.gcr.io/busybox:latest",
			objs:     container("k8s.gcr.io/busybox:1.2.3"),
			expected: nil,
		},
		{
			name:     "pod with container image - busybox:latest",
			objs:     container("busybox:latest"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "pod with container image - k8s.gcr.io/busybox",
			objs:     container("k8s.gcr.io/busybox"),
			expected: nil,
		},
		{
			name:     "pod with container image - busybox",
			objs:     container("busybox"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "pod with container image - test:5000/repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     container("test:5000/repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: nil,
		},
		{
			name:     "pod with container image - repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     container("repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "pod with container image - test:5000/repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     container("test:5000/repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: nil,
		},
		{
			name:     "pod with container image - repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     container("repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "pod with container image - k8s.gcr.io/busybox:latest",
			objs:     initContainer("k8s.gcr.io/busybox:1.2.3"),
			expected: nil,
		},
		{
			name:     "pod with container image - busybox:latest",
			objs:     initContainer("busybox:latest"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "pod with container image - k8s.gcr.io/busybox",
			objs:     initContainer("k8s.gcr.io/busybox"),
			expected: nil,
		},
		{
			name:     "pod with container image - busybox",
			objs:     initContainer("busybox"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "pod with container image - test:5000/repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     initContainer("test:5000/repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: nil,
		},
		{
			name:     "pod with container image - repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     initContainer("repo/image@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "pod with container image - test:5000/repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     initContainer("test:5000/repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: nil,
		},
		{
			name:     "pod with container image - repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			objs:     initContainer("repo/image:ignore-tag@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			expected: issues(severity, message, check.Pod, name),
		},
	}

	fullyQualifiedImageCheck := fullyQualifiedImageCheck{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := fullyQualifiedImageCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}

func TestMalformedImageError(t *testing.T) {
	const message = "Malformed image name for container 'bar'"
	const severity = check.Warning
	const name = "fully-qualified-image"

	tests := []struct {
		name     string
		objs     *check.Objects
		expected []check.Diagnostic
	}{
		{
			name:     "container with image : test:5000/repo/image@sha256:digest",
			objs:     container("test:5000/repo/image@sha256:digest"),
			expected: issues(severity, message, check.Pod, name),
		},
		{
			name:     "init container with image : test:5000/repo/image@sha256:digest",
			objs:     initContainer("test:5000/repo/image@sha256:digest"),
			expected: issues(severity, message, check.Pod, name),
		},
	}
	fullyQualifiedImageCheck := fullyQualifiedImageCheck{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, _, err := fullyQualifiedImageCheck.Run(test.objs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expected, d)
		})
	}
}
