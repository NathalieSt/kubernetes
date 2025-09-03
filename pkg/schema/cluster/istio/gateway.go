package istio

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

/*
import k8s.apimachinery.pkg.apis.meta.v1 as meta
import ..shared

schema Server:
    port: shared.IstioPort
    hosts: [str]

schema GatewaySelector:
    istio: str

schema GatewaySpec:
    servers: [Server]
    selector: GatewaySelector

schema Gateway:
    apiVersion = "networking.istio.io/v1"
    kind = "Gateway"
    metadata: meta.ObjectMeta
    spec: GatewaySpec


*/

type GatewayServer struct {
	Port  IstioPort `yaml:",omitempty"`
	Hosts []string  `yaml:",omitempty"`
}

type GatewaySelector struct {
	Istio string `yaml:",omitempty"`
}

type GatewaySpec struct {
	Servers  []GatewayServer `yaml:",omitempty"`
	Selector GatewaySelector `yaml:",omitempty"`
}

type Gateway struct {
	shared.CommonK8sResourceWithSpec[GatewaySpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewGateway(meta meta.ObjectMeta, spec GatewaySpec) Gateway {
	return Gateway{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[GatewaySpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "networking.istio.io/v1",
				Kind:       "Gateway",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
