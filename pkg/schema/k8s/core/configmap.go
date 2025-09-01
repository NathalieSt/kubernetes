package core

import (
	"kubernetes/pkg/schema/k8s/meta"
)

type ConfigMap struct {
	ApiVersion string `yaml:"apiVersion,omitempty" validate:"required"`
	Kind       string `yaml:"kind,omitempty" validate:"required"`
	Metadata   meta.ObjectMeta
	Data       map[string]string
}
