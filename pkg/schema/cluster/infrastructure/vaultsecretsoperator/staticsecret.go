package vaultsecretsoperator

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Destination struct {
	Create      bool              `yaml:"create"`
	Name        string            `yaml:"name"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
}
type StaticSecretSpec struct {
	AuthRef      string      `yaml:"AuthRef"`
	Mount        string      `yaml:"mount"`
	Type         string      `yaml:"type"`
	Path         string      `yaml:"path"`
	RefreshAfter string      `yaml:"refreshAfter"`
	Destination  Destination `yaml:"destination"`
}

type StaticSecret struct {
	shared.CommonK8sResourceWithSpec[StaticSecretSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewStaticSecret(meta meta.ObjectMeta, spec StaticSecretSpec) StaticSecret {
	return StaticSecret{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[StaticSecretSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "secrets.hashicorp.com/v1beta1",
				Kind:       "StaticSecret",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
