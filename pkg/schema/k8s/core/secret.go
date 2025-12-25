package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Secret struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	Data                     map[string][]byte `yaml:"data,omitempty"`
	Immutable                bool              `yaml:"immutable,omitempty"`
	StringData               map[string]string `yaml:"stringData,omitempty"`
	Type                     string            `yaml:"type,omitempty"`
}

type SecretConfig struct {
	Data       map[string][]byte `yaml:"data,omitempty"`
	Immutable  bool              `yaml:"immutable,omitempty"`
	StringData map[string]string `yaml:"stringData,omitempty"`
	Type       string            `yaml:"type,omitempty"`
}

func NewSecret(meta meta.ObjectMeta, config SecretConfig) Secret {
	return Secret{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "v1",
			Kind:       "Secret",
			Metadata:   meta,
		},
		Data:       config.Data,
		Immutable:  config.Immutable,
		StringData: config.StringData,
		Type:       config.Type,
	}
}
