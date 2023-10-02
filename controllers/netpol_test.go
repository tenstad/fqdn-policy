package controllers

import (
	"github.com/GoogleCloudPlatform/gke-fqdnnetworkpolicies-golang/api/v1alpha3"
	"github.com/stretchr/testify/assert"
	networking "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

type TestCase struct {
	currentFqdn       *v1alpha3.FQDNNetworkPolicy
	foundEgressRules  []networking.NetworkPolicyEgressRule
	foundIngressRules []networking.NetworkPolicyIngressRule
}

func TestNetpol(t *testing.T) {
	tests := []struct {
		name          string
		testCase      TestCase
		currentNetpol Netpol
		shouldUpdate  bool
	}{
		{
			name:          "should not update network policy",
			testCase:      newTestCase("10.10", "10.10"),
			currentNetpol: netPool("10.10", "10.10"),
		}, {
			name:          "found egress changed, should update network policy",
			testCase:      newTestCase("10.10", "10.12"),
			currentNetpol: netPool("10.10", "10.11"),
			shouldUpdate:  true,
		}, {
			name:          "current netpol egress differs, should update network policy",
			testCase:      newTestCase("10.10", "10.11"),
			currentNetpol: netPool("10.10", "10.12"),
			shouldUpdate:  true,
		}, {
			name:          "found ingress differs, should update network policy",
			testCase:      newTestCase("10.11", "10.10"),
			currentNetpol: netPool("10.10", "10.10"),
			shouldUpdate:  true,
		}, {
			name:          "current netpol ingress differs, should update network policy",
			testCase:      newTestCase("10.11", "10.10"),
			currentNetpol: netPool("10.12", "10.10"),
			shouldUpdate:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// check of network policy
			tt.currentNetpol.UpdateMetadata(tt.testCase.currentFqdn)
			assert.Equal(t, tt.testCase.currentFqdn.Name, tt.currentNetpol.Name)
			assert.Equal(t, tt.testCase.currentFqdn.Namespace, tt.currentNetpol.Namespace)
			assert.Equal(t, tt.testCase.currentFqdn.Annotations, tt.currentNetpol.Annotations)
			assert.Equal(t, tt.testCase.currentFqdn.Spec.PodSelector, tt.currentNetpol.Spec.PodSelector)
			assert.Equal(t, tt.testCase.currentFqdn.Spec.PolicyTypes, tt.currentNetpol.Spec.PolicyTypes)

			if tt.shouldUpdate {
				tt.currentNetpol.UpdateEgress(tt.testCase.foundEgressRules)
				tt.currentNetpol.UpdateIngress(tt.testCase.foundIngressRules)
				assert.Equal(t, tt.testCase.foundEgressRules, tt.currentNetpol.Spec.Egress)
				assert.Equal(t, tt.testCase.foundIngressRules, tt.currentNetpol.Spec.Ingress)
			}
		})
	}
}

func newTestCase(ingressCidr, egressCidr string) TestCase {
	return TestCase{
		currentFqdn: &v1alpha3.FQDNNetworkPolicy{
			ObjectMeta: v1.ObjectMeta{
				Name:      "my-fqdn",
				Namespace: "default",
				Annotations: map[string]string{
					ownerAnnotation: "my-fqdn",
				},
			},
			Spec: v1alpha3.FQDNNetworkPolicySpec{
				PolicyTypes: []networking.PolicyType{
					networking.PolicyTypeEgress,
				},
				Egress: []v1alpha3.FQDNNetworkPolicyEgressRule{
					{
						To: []v1alpha3.FQDNNetworkPolicyPeer{
							{
								FQDNs: []string{"my-fqdn"},
							},
						},
					},
				},
			},
		},
		foundEgressRules:  egress(egressCidr),
		foundIngressRules: ingress(ingressCidr),
	}
}

func ingress(cidr string) []networking.NetworkPolicyIngressRule {
	return []networking.NetworkPolicyIngressRule{
		{
			From: []networking.NetworkPolicyPeer{
				{
					IPBlock: &networking.IPBlock{
						CIDR: cidr,
					},
				},
			},
		},
	}
}

func egress(cidr string) []networking.NetworkPolicyEgressRule {
	return []networking.NetworkPolicyEgressRule{
		{
			To: []networking.NetworkPolicyPeer{
				{
					IPBlock: &networking.IPBlock{
						CIDR: cidr,
					},
				},
			},
		},
	}
}

func netPool(ingressCidr, egressCidr string) Netpol {
	return Netpol{
		NetworkPolicy: &networking.NetworkPolicy{
			Spec: networking.NetworkPolicySpec{
				PolicyTypes: []networking.PolicyType{
					networking.PolicyTypeEgress,
				},
				Egress:  egress(egressCidr),
				Ingress: ingress(ingressCidr),
			},
		},
	}
}
