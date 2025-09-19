package istio

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ServiceEntryPorts struct {
	Number   int64  `yaml:"number,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
}
type ServiceEntrySpec struct {
	Hosts      []string            `yaml:"hosts,omitempty"`
	Ports      []ServiceEntryPorts `yaml:"ports,omitempty"`
	Location   string              `yaml:"location,omitempty"`
	Resolution string              `yaml:"resolution,omitempty"`
}

type ServiceEntry struct {
	shared.CommonK8sResourceWithSpec[ServiceEntrySpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewServiceEntry(meta meta.ObjectMeta, spec ServiceEntrySpec) ServiceEntry {
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
