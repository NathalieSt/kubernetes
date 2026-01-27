package vaultsecretsoperator

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type SecretType = string

const (
	Opaque    SecretType = "Opaque"
	BasicAuth SecretType = "kubernetes.io/basic-auth"
)

type Template struct {
	Name string `yaml:"name,omitempty"`
	Text string `yaml:"text,omitempty"`
}

type TemplateRef struct {
	Name        string `yaml:"name,omitempty"`
	KeyOverride string `yaml:"keyOverride,omitempty"`
}

type TransformationRef struct {
	Namespace      string        `yaml:"namespace,omitempty"`
	Name           string        `yaml:"name,omitempty"`
	TemplateRefs   []TemplateRef `yaml:"templateRefs,omitempty"`
	IgnoreIncludes bool          `yaml:"ignoreIncludes,omitempty"`
	IgnoreExcludes bool          `yaml:"ignoreExcludes,omitempty"`
}

type Transformation struct {
	Templates          map[string]Template `yaml:"templates,omitempty"`
	TransformationRefs []TransformationRef `yaml:"transformationRefs,omitempty"`
	Includes           []string            `yaml:"includes,omitempty"`
	Excludes           []string            `yaml:"excludes,omitempty"`
	ExcludeRaw         bool                `yaml:"excludeRaw,omitempty"`
}

type Destination struct {
	Create         bool              `yaml:"create,omitempty"`
	Name           string            `yaml:"name,omitempty"`
	Labels         map[string]string `yaml:"labels,omitempty"`
	Annotations    map[string]string `yaml:"annotations,omitempty"`
	Type           SecretType        `yaml:"type,omitempty"`
	Transformation Transformation    `yaml:"transformation,omitempty"`
}

type StaticSecretSpec struct {
	AuthRef      string      `yaml:"vaultAuthRef,omitempty"`
	Mount        string      `yaml:"mount,omitempty"`
	Type         string      `yaml:"type,omitempty"`
	Path         string      `yaml:"path,omitempty"`
	RefreshAfter string      `yaml:"refreshAfter,omitempty"`
	Destination  Destination `yaml:"destination,omitempty"`
}

type StaticSecret struct {
	shared.CommonK8sResourceWithSpec[StaticSecretSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewStaticSecret(meta meta.ObjectMeta, spec StaticSecretSpec) StaticSecret {
	return StaticSecret{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[StaticSecretSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "secrets.hashicorp.com/v1beta1",
				Kind:       "VaultStaticSecret",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
