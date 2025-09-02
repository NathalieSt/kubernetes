package core

import (
	"kubernetes/pkg/schema/k8s/meta"
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
	ApiVersion string `yaml:",apiVersion"`
	Kind       string
	Metadata   meta.ObjectMeta
	Spec       ServiceSpec
}

func NewService(meta meta.ObjectMeta, spec ServiceSpec) Service {
	return Service{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   meta,
		Spec:       spec,
	}
}
