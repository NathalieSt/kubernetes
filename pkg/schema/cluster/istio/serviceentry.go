package istio

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ServiceEntryPorts struct {
	Number   int    `yaml:"number"`
	Name     string `yaml:"name"`
	Protocol string `yaml:"protocol"`
}
type ServiceEntrySpec struct {
	Hosts      []string            `yaml:"hosts"`
	Ports      []ServiceEntryPorts `yaml:"ports"`
	Location   string              `yaml:"location"`
	Resolution string              `yaml:"resolution"`
}

type ServiceEntry struct {
	shared.CommonK8sResourceWithSpec[ServiceEntrySpec] `yaml:",omitempty,inline" validate:"required"`
	Spec                                               ServiceEntrySpec
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
