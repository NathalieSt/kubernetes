package flux

import "kubernetes/pkg/schema/k8s/meta"

type Kustomization struct {
	ApiVersion string `yaml:"apiVersion," validate:"required"`
	Kind       string `yaml:"kind," validate:"required"`
	Metadata   meta.ObjectMeta
	Resources  []string
}

func NewKustomization(meta meta.ObjectMeta, resources []string) Kustomization {
	return Kustomization{
		ApiVersion: "source.toolkit.fluxcd.io/v1",
		Kind:       "HelmRepository",
		Metadata:   meta,
		Resources:  resources,
	}
}
