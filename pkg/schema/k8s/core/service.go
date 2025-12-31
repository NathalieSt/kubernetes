package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Protocol = string

const (
	SCTP Protocol = "SCTP"
	TCP  Protocol = "TCP"
	UDP  Protocol = "UDP"
)

type ServicePort struct {
	Name       string   `yaml:"name,omitempty"`
	Port       int64    `yaml:"port,omitempty"`
	TargetPort int64    `yaml:"targetPort,omitempty"`
	Protocol   Protocol `yaml:"protocol,omitempty"`
}

type ServiceSpec struct {
	Selector map[string]string `yaml:"selector,omitempty"`
	Ports    []ServicePort     `yaml:"ports,omitempty"`
	Type     string            `yaml:"type,omitempty"`
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
