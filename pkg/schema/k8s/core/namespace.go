package core

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type Namespace struct {
	ApiVersion string `yaml:"apiVersion,omitempty" validate:"required"`
	Kind       string `yaml:"kind,omitempty" validate:"required"`
	Metadata   meta.ObjectMeta
}

func NewNamespace(metadata meta.ObjectMeta) Namespace {
	return Namespace{
		ApiVersion: "v1",
		Kind:       "Namespace",
		Metadata:   metadata,
	}
}
