package kustomize

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Kustomization struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	Resources                []string
}

func NewKustomization(meta meta.ObjectMeta, resources []string) Kustomization {
	return Kustomization{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "kustomize.config.k8s.io/v1beta1",
			Kind:       "Kustomization",
			Metadata:   meta,
		},
		Resources: resources,
	}
}
