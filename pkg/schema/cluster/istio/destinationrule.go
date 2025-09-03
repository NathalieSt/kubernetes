package istio

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type MTLSModes string

const (
	DESTINATION_DISABLE MTLSModes = "DISABLE"
	SIMPLE              MTLSModes = "SIMPLE"
	MUTUAL              MTLSModes = "MUTUAL"
	ISTIO_MUTUAL        MTLSModes = "ISTIO_MUTUAL"
)

type DestinationRuleTrafficPolicyTLS struct {
	Mode MTLSModes
}

type DestinationRuleTrafficPolicy struct {
	TLS DestinationRuleTrafficPolicyTLS
}

type DestinationRuleSpec struct {
	Host          string
	TrafficPolicy DestinationRuleTrafficPolicy
}

type DestinationRule struct {
	shared.CommonK8sResourceWithSpec[DestinationRuleSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewDestinationRule(meta meta.ObjectMeta, spec DestinationRuleSpec) DestinationRule {
	return DestinationRule{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[DestinationRuleSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "networking.istio.io/v1",
				Kind:       "DestinationRule",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
