package core

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type ServicePort struct {
	Name       string
	Port       int
	TargetPort int
}

type ServiceSpec struct {
	Selector map[string]string
	Ports    []ServicePort
	Type     string
}

type Service struct {
	ApiVersion string
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
