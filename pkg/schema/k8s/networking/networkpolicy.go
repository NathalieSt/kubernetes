package networking

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type IPBlock struct {
	CIDR   string   `yaml:"cidr,omitempty"`
	Except []string `yaml:"except,omitempty"`
}

type NetworkPolicyPeer struct {
	IpBlock           IPBlock            `yaml:"ipBlock,omitempty"`
	NamespaceSelector meta.LabelSelector `yaml:"namespaceSelector,omitempty"`
	PodSelector       meta.LabelSelector `yaml:"podSelector,"`
}

type PortProtocol string

const (
	SCTP PortProtocol = "SCTP"
	TCP  PortProtocol = "TCP"
	UDP  PortProtocol = "UDP"
)

type NetworkPolicyPort struct {
	Port     int32        `yaml:"port,omitempty"`
	EndPort  int32        `yaml:"endPort,omitempty"`
	Protocol PortProtocol `yaml:"protocol,omitempty"`
}

type NetworkPolicyEgressRule struct {
	To    []NetworkPolicyPeer `yaml:"to,"`
	Ports []NetworkPolicyPort `yaml:"ports,omitempty"`
}

type NetworkPolicyIngressRule struct {
	From  []NetworkPolicyPeer `yaml:"from,"`
	Ports []NetworkPolicyPort `yaml:"ports,omitempty"`
}

type NetworkPolicyType string

const (
	Ingress NetworkPolicyType = "Ingress"
	Egress  NetworkPolicyType = "Egress"
)

type NetworkPolicySpec struct {
	PodSelector meta.LabelSelector         `yaml:"podSelector,omitempty"`
	PolicyTypes []NetworkPolicyType        `yaml:"policyTypes,omitempty"`
	Ingress     []NetworkPolicyIngressRule `yaml:"ingress,omitempty"`
	Egress      []NetworkPolicyEgressRule  `yaml:"egress,omitempty"`
}

type NetworkPolicy struct {
	shared.CommonK8sResourceWithSpec[NetworkPolicySpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewNetworkPolicy(meta meta.ObjectMeta, spec NetworkPolicySpec) NetworkPolicy {
	return NetworkPolicy{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[NetworkPolicySpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "networking.k8s.io/v1",
				Kind:       "NetworkPolicy",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
