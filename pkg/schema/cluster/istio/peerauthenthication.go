package istio

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

/*

schema Mtls:
    # FIXME: theres a third option, like "none" or something like that
    mode: "STRICT" | "PERMISSIVE"

schema Selector:
    matchLabels: {str:str}

schema PeerAuthenthicationSpec:
    selector?: Selector
    mtls: Mtls

schema PeerAuthenthication:
    apiVersion = "security.istio.io/v1"
    kind = "PeerAuthentication"
    metadata: meta.ObjectMeta
    spec: PeerAuthenthicationSpec


*/

type PeerAuthenthicationmTLSMode string

const (
	UNSET        PeerAuthenthicationmTLSMode = "UNSET"
	PEER_DISABLE PeerAuthenthicationmTLSMode = "DISABLE"
	PERMISSIVE   PeerAuthenthicationmTLSMode = "PERMISSIVE"
	STRICT       PeerAuthenthicationmTLSMode = "STRICT"
)

type PeerAuthenthicationmTLS struct {
	Mode PeerAuthenthicationmTLSMode
}

type PeerAuthenthicationSelector struct {
	meta.LabelSelector `yaml:",omitempty,inline"`
}

type PeerAuthenthicationSpec struct {
	Selector PeerAuthenthicationSelector
	MTLS     PeerAuthenthicationmTLS
}

type PeerAuthenthication struct {
	shared.CommonK8sResourceWithSpec[PeerAuthenthicationSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewPeerAuthenthication(meta meta.ObjectMeta, spec PeerAuthenthicationSpec) PeerAuthenthication {
	return PeerAuthenthication{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[PeerAuthenthicationSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "security.istio.io/v1",
				Kind:       "PeerAuthentication",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
