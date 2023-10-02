package controllers

import (
	"fmt"
	"github.com/GoogleCloudPlatform/gke-fqdnnetworkpolicies-golang/api/v1alpha3"
	"github.com/mitchellh/hashstructure/v2"
	networking "k8s.io/api/networking/v1"
)

const (
	ownerAnnotation = "fqdnnetworkpolicies.networking.gke.io/owned-by"
)

type Netpol struct {
	*networking.NetworkPolicy
}

func (n *Netpol) UpdateMetadata(fqdn *v1alpha3.FQDNNetworkPolicy) {
	n.Name = fqdn.Name
	n.Namespace = fqdn.Namespace
	if n.Annotations == nil {
		n.Annotations = make(map[string]string)
	}
	n.Annotations[ownerAnnotation] = fqdn.Name
	n.Spec.PodSelector = fqdn.Spec.PodSelector
	n.Spec.PolicyTypes = fqdn.Spec.PolicyTypes
}

func (n *Netpol) UpdateEgress(fqdnEgress []networking.NetworkPolicyEgressRule) {
	n.Spec.Egress = fqdnEgress
}

func (n *Netpol) UpdateIngress(fqdnIngress []networking.NetworkPolicyIngressRule) {
	n.Spec.Ingress = fqdnIngress
}

func (n *Netpol) EgressRulesEquals(egress []networking.NetworkPolicyEgressRule) (bool, error) {
	netpolHash, err := hashstructure.Hash(n.NetworkPolicy.Spec.Egress, hashstructure.FormatV2, nil)
	if err != nil {
		return false, err
	}

	currentHash, err := hashstructure.Hash(egress, hashstructure.FormatV2, nil)
	if err != nil {
		return false, err
	}
	return fmt.Sprintf("%x", netpolHash) == fmt.Sprintf("%x", currentHash), nil
}

func (n *Netpol) IngressRulesEquals(ingress []networking.NetworkPolicyIngressRule) (bool, error) {
	netpolHash, err := hashstructure.Hash(n.NetworkPolicy.Spec.Ingress, hashstructure.FormatV2, nil)
	if err != nil {
		return false, err
	}

	currentHash, err := hashstructure.Hash(ingress, hashstructure.FormatV2, nil)
	if err != nil {
		return false, err
	}
	return fmt.Sprintf("%x", netpolHash) == fmt.Sprintf("%x", currentHash), nil
}
