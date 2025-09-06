package istio

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type VirtualServicePortSelector struct {
	Number int
}

type VirtualServiceDestination struct {
	Host string
	Port VirtualServicePortSelector
}

type VirtualServiceSetProtocolType string

const (
	VirtualServiceHTTPProto  VirtualServiceSetProtocolType = "http"
	VirtualServiceHTTPSProto VirtualServiceSetProtocolType = "https"
)

type VirtualServiceSet struct {
	XForwardedProto VirtualServiceSetProtocolType `yaml:"X-Forwarded-Proto,omitempty"`
}

type VirtualServiceRequest struct {
	Set VirtualServiceSet
}

type VirtualServiceHeaders struct {
	Request VirtualServiceRequest
}

type VirtualServiceRoute struct {
	Destination VirtualServiceDestination
	Headers     VirtualServiceHeaders
}

type VirtualServiceHTTP struct {
	Route []VirtualServiceRoute
}

type VirtualServiceSpec struct {
	Hosts    []string
	Gateways []string
	HTTP     []VirtualServiceHTTP
}

type VirtualService struct {
	shared.CommonK8sResourceWithSpec[VirtualServiceSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewVirtualService(meta meta.ObjectMeta, spec VirtualServiceSpec) VirtualService {
	return VirtualService{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[VirtualServiceSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "networking.istio.io/v1",
				Kind:       "VirtualService",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
