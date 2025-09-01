package core

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type ServicePort struct {
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
}
