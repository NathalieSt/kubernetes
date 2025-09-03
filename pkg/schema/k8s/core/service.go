package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ServicePort struct {
	Name       string `yaml:",omitempty"`
	Port       int    `yaml:",omitempty"`
	TargetPort int    `yaml:",omitempty"`
}

type ServiceSpec struct {
	Selector map[string]string `yaml:",omitempty"`
	Ports    []ServicePort     `yaml:",omitempty"`
	Type     string            `yaml:",omitempty"`
}

type Service struct {
	shared.CommonK8sResourceWithSpec[ServiceSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewService(meta meta.ObjectMeta, spec ServiceSpec) Service {
	return Service{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[ServiceSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "v1",
				Kind:       "Service",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
