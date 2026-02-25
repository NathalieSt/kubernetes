package shared

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/k8s/networking"
)

var NetbirdDomainBase = "netbird.nathalie-stiefsohn.eu"

var ExternalIngressNetworkPolicyPeer = networking.NetworkPolicyPeer{
	PodSelector: meta.LabelSelector{
		MatchLabels: map[string]string{
			"app.kubernetes.io/name": "netbird-router",
		},
	},
	NamespaceSelector: meta.LabelSelector{
		MatchLabels: map[string]string{
			"kubernetes.io/metadata.name": "netbird",
		},
	},
}
