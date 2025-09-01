package k8s

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type Kustomization struct {
	ApiVersion string `yaml:"apiVersion," validate:"required"`
	Kind       string `yaml:"kind," validate:"required"`
	Metadata   meta.ObjectMeta
	Resources  []string
}
