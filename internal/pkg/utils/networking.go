package utils

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
)

func GetCoreDNSEgressRule() networking.NetworkPolicyEgressRule {
	rule := networking.NetworkPolicyEgressRule{
		Ports: []networking.NetworkPolicyPort{
			{
				Port:     53,
				Protocol: networking.UDP,
			},
			{
				Port:     53,
				Protocol: networking.TCP,
			},
		},
		To: []networking.NetworkPolicyPeer{
			{
				NamespaceSelector: meta.LabelSelector{
					MatchLabels: map[string]string{
						"kubernetes.io/metadata.name": "kube-system",
					},
				},
				PodSelector: meta.LabelSelector{
					MatchLabels: map[string]string{
						"k8s-app": "kube-dns",
					},
				},
			},
		},
	}
	return rule
}
