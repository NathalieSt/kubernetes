package istio

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

/*
schema ServiceEntrySpec:
    hosts: [str]
    ports: [shared.IstioPort]
    location: "MESH_EXTERNAL" | "MESH_INTERNAL"
    resolution: "DNS" | "STATIC"

schema ServiceEntry:
    apiVersion = "networking.istio.io/v1alpha3"
    kind = "ServiceEntry"
    metadata: meta.ObjectMeta
    spec: ServiceEntrySpec

*/

type ServiceEntrySpec struct {
}

type ServiceEntry struct {
	shared.CommonK8sResourceWithSpec[ServiceEntrySpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewPeerServiceEntry(meta meta.ObjectMeta, spec ServiceEntrySpec) ServiceEntry {
	return ServiceEntry{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[ServiceEntrySpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "networking.istio.io/v1alpha3",
				Kind:       "ServiceEntry",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
