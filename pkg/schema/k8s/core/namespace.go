package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Namespace struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
}

func NewNamespace(metadata meta.ObjectMeta) Namespace {
	return Namespace{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "v1",
			Kind:       "Namespace",
			Metadata:   metadata,
		},
	}
}
