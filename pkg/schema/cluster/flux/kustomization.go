package flux

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
			ApiVersion: "source.toolkit.fluxcd.io/v1",
			Kind:       "HelmRepository",
			Metadata:   meta,
		},
		Resources: resources,
	}
}
