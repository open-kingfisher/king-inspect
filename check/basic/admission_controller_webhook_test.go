package basic

import (
	"testing"

	"github.com/open-kingfisher/king-inspect/check"
	"github.com/stretchr/testify/assert"
	ar "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWebhookCheckMeta(t *testing.T) {
	webhookCheck := webhookCheck{}
	assert.Equal(t, "admission-controller-webhook", webhookCheck.Name())
	assert.Equal(t, []string{"basic"}, webhookCheck.Groups())
	assert.NotEmpty(t, webhookCheck.Description())
}

func TestWebhookCheckRegistration(t *testing.T) {
	webhookCheck := &webhookCheck{}
	check, err := check.Get("admission-controller-webhook")
	assert.NoError(t, err)
	assert.Equal(t, check, webhookCheck)
}

func TestWebHookRun(t *testing.T) {
	emptyNamespaceList := &corev1.NamespaceList{
		Items: []corev1.Namespace{},
	}
	emptyServiceList := &corev1.ServiceList{
		Items: []corev1.Service{},
	}

	baseMWC := ar.MutatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{Kind: "MutatingWebhookConfiguration", APIVersion: "v1beta1"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwc_foo",
		},
		Webhooks: []ar.MutatingWebhook{},
	}
	baseMW := ar.MutatingWebhook{
		Name: "mw_foo",
	}

	baseVWC := ar.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{Kind: "ValidatingWebhookConfiguration", APIVersion: "v1beta1"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwc_foo",
		},
		Webhooks: []ar.ValidatingWebhook{},
	}
	baseVW := ar.ValidatingWebhook{
		Name: "vw_foo",
	}

	tests := []struct {
		name     string
		objs     *check.Objects
		expected []check.Diagnostic
	}{
		{
			name: "no webhook configuration",
			objs: &check.Objects{
				MutatingWebhookConfigurations:   &ar.MutatingWebhookConfigurationList{},
				ValidatingWebhookConfigurations: &ar.ValidatingWebhookConfigurationList{},
				SystemNamespace:                 &corev1.Namespace{},
			},
			expected: nil,
		},
		{
			name: "direct url webhooks",
			objs: &check.Objects{
				Namespaces: emptyNamespaceList,
				MutatingWebhookConfigurations: &ar.MutatingWebhookConfigurationList{
					Items: []ar.MutatingWebhookConfiguration{
						func() ar.MutatingWebhookConfiguration {
							mwc := baseMWC
							mw := baseMW
							mw.ClientConfig = ar.WebhookClientConfig{
								URL: strPtr("http://webhook.com"),
							}
							mwc.Webhooks = append(mwc.Webhooks, mw)
							return mwc
						}(),
					},
				},
				ValidatingWebhookConfigurations: &ar.ValidatingWebhookConfigurationList{
					Items: []ar.ValidatingWebhookConfiguration{
						func() ar.ValidatingWebhookConfiguration {
							vwc := baseVWC
							vw := baseVW
							vw.ClientConfig = ar.WebhookClientConfig{
								URL: strPtr("http://webhook.com"),
							}
							vwc.Webhooks = append(vwc.Webhooks, vw)
							return vwc
						}(),
					},
				},
				SystemNamespace: &corev1.Namespace{},
			},
			expected: nil,
		},
		{
			name: "namespace does not exist",
			objs: &check.Objects{
				Namespaces: emptyNamespaceList,
				Services: &corev1.ServiceList{
					Items: []corev1.Service{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "service",
							},
						},
					},
				},
				MutatingWebhookConfigurations: &ar.MutatingWebhookConfigurationList{
					Items: []ar.MutatingWebhookConfiguration{
						func() ar.MutatingWebhookConfiguration {
							mwc := baseMWC
							mw := baseMW
							mw.ClientConfig = ar.WebhookClientConfig{
								Service: &ar.ServiceReference{
									Namespace: "missing",
									Name:      "service",
								},
							}
							mwc.Webhooks = append(mwc.Webhooks, mw)
							return mwc
						}(),
					},
				},
				ValidatingWebhookConfigurations: &ar.ValidatingWebhookConfigurationList{
					Items: []ar.ValidatingWebhookConfiguration{
						func() ar.ValidatingWebhookConfiguration {
							vwc := baseVWC
							vw := baseVW
							vw.ClientConfig = ar.WebhookClientConfig{
								Service: &ar.ServiceReference{
									Namespace: "missing",
									Name:      "service",
								},
							}
							vwc.Webhooks = append(vwc.Webhooks, vw)
							return vwc
						}(),
					},
				},
				SystemNamespace: &corev1.Namespace{},
			},
			expected: []check.Diagnostic{
				{
					Severity: check.Error,
					Message:  "Validating webhook vw_foo is configured against a service in a namespace that does not exist.",
					Kind:     check.ValidatingWebhookConfiguration,
				},
				{
					Severity: check.Error,
					Message:  "Mutating webhook mw_foo is configured against a service in a namespace that does not exist.",
					Kind:     check.MutatingWebhookConfiguration,
				},
			},
		},
		{
			name: "service does not exist",
			objs: &check.Objects{
				Namespaces: &corev1.NamespaceList{
					Items: []corev1.Namespace{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "webhook",
							},
						},
					},
				},
				Services: emptyServiceList,
				MutatingWebhookConfigurations: &ar.MutatingWebhookConfigurationList{
					Items: []ar.MutatingWebhookConfiguration{
						func() ar.MutatingWebhookConfiguration {
							mwc := baseMWC
							mw := baseMW
							mw.ClientConfig = ar.WebhookClientConfig{
								Service: &ar.ServiceReference{
									Namespace: "webhook",
									Name:      "service",
								},
							}
							mwc.Webhooks = append(mwc.Webhooks, mw)
							return mwc
						}(),
					},
				},
				ValidatingWebhookConfigurations: &ar.ValidatingWebhookConfigurationList{
					Items: []ar.ValidatingWebhookConfiguration{
						func() ar.ValidatingWebhookConfiguration {
							vwc := baseVWC
							vw := baseVW
							vw.ClientConfig = ar.WebhookClientConfig{
								Service: &ar.ServiceReference{
									Namespace: "webhook",
									Name:      "service",
								},
							}
							vwc.Webhooks = append(vwc.Webhooks, vw)
							return vwc
						}(),
					},
				},
				SystemNamespace: &corev1.Namespace{},
			},
			expected: []check.Diagnostic{
				{
					Severity: check.Error,
					Message:  "Validating webhook vw_foo is configured against a service that does not exist.",
					Kind:     check.ValidatingWebhookConfiguration,
				},
				{
					Severity: check.Error,
					Message:  "Mutating webhook mw_foo is configured against a service that does not exist.",
					Kind:     check.MutatingWebhookConfiguration,
				},
			},
		},
	}

	webhookCheck := webhookCheck{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diagnostics, _, err := webhookCheck.Run(test.objs)
			assert.NoError(t, err)

			// skip checking object and owner since for this checker it just uses the object being checked.
			var strippedDiagnostics []check.Diagnostic
			for _, d := range diagnostics {
				d.Object = nil
				d.Owners = nil
				strippedDiagnostics = append(strippedDiagnostics, d)
			}

			assert.ElementsMatch(t, test.expected, strippedDiagnostics)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
