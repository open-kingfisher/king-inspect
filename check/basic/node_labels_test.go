package basic

import (
	"testing"

	"github.com/open-kingfisher/king-inspect/check"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodeLabels(t *testing.T) {
	tests := []struct {
		name                string
		nodeLabels          map[string]string
		expectedDiagnostics []check.Diagnostic
	}{
		{
			name:                "no labels",
			nodeLabels:          nil,
			expectedDiagnostics: nil,
		},
		{
			name: "only doks labels",
			nodeLabels: map[string]string{
				"doks.digitalocean.com/foo": "bar",
				"doks.digitalocean.com/baz": "xyzzy",
			},
			expectedDiagnostics: nil,
		},
		{
			name: "only built-in labels",
			nodeLabels: map[string]string{
				"checkrnetes.io/hostname":                   "a-hostname",
				"beta.checkrnetes.io/os":                    "linux",
				"failure-domain.beta.checkrnetes.io/region": "tor1",
			},
			expectedDiagnostics: nil,
		},
		{
			name: "only region label",
			nodeLabels: map[string]string{
				"region": "tor1",
			},
			expectedDiagnostics: nil,
		},
		{
			name: "custom labels",
			nodeLabels: map[string]string{
				"doks.digitalocean.com/foo":                 "bar",
				"doks.digitalocean.com/baz":                 "xyzzy",
				"checkrnetes.io/hostname":                   "a-hostname",
				"example.com/custom-label":                  "bad",
				"example.com/another-label":                 "real-bad",
				"beta.checkrnetes.io/os":                    "linux",
				"failure-domain.beta.checkrnetes.io/region": "tor1",
				"region": "tor1",
			},
			expectedDiagnostics: []check.Diagnostic{{
				Severity: check.Warning,
				Message:  "Custom node labels will be lost if node is replaced or upgraded.",
				Kind:     check.Node,
				Object: &metav1.ObjectMeta{
					Labels: map[string]string{
						"doks.digitalocean.com/foo":                 "bar",
						"doks.digitalocean.com/baz":                 "xyzzy",
						"checkrnetes.io/hostname":                   "a-hostname",
						"example.com/custom-label":                  "bad",
						"example.com/another-label":                 "real-bad",
						"beta.checkrnetes.io/os":                    "linux",
						"failure-domain.beta.checkrnetes.io/region": "tor1",
						"region": "tor1",
					},
				},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			objects := &check.Objects{
				Nodes: &corev1.NodeList{
					Items: []corev1.Node{{
						ObjectMeta: metav1.ObjectMeta{
							Labels: test.nodeLabels,
						},
					}},
				},
			}

			check := &nodeLabelsTaintsCheck{}

			ds, _, err := check.Run(objects)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expectedDiagnostics, ds)
		})
	}
}

func TestNodeTaints(t *testing.T) {
	tests := []struct {
		name                string
		taints              []corev1.Taint
		expectedDiagnostics []check.Diagnostic
	}{
		{
			name:                "no taints",
			taints:              nil,
			expectedDiagnostics: nil,
		},
		{
			name: "custom taints",
			taints: []corev1.Taint{{
				Key:    "example.com/my-taint",
				Value:  "foo",
				Effect: corev1.TaintEffectNoSchedule,
			}},
			expectedDiagnostics: []check.Diagnostic{{
				Severity: check.Warning,
				Message:  "Custom node taints will be lost if node is replaced or upgraded.",
				Kind:     check.Node,
				Object:   &metav1.ObjectMeta{},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			objects := &check.Objects{
				Nodes: &corev1.NodeList{
					Items: []corev1.Node{{
						Spec: corev1.NodeSpec{
							Taints: test.taints,
						},
					}},
				},
			}

			check := &nodeLabelsTaintsCheck{}

			ds, _, err := check.Run(objects)
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.expectedDiagnostics, ds)
		})
	}
}
