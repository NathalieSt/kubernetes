package vaultsecretsoperator

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type AuthGlobalRef struct {
	AllowDefault bool   `yaml:"allowDefault,omitempty"`
	Namespace    string `yaml:"namespace,omitempty"`
}
type AuthSpec struct {
	Kubernetes         Kubernetes    `yaml:"kubernetes,omitempty"`
	VaultAuthGlobalRef AuthGlobalRef `yaml:"vaultAuthGlobalRef,omitempty"`
}

type Auth struct {
	shared.CommonK8sResourceWithSpec[AuthSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewAuth(meta meta.ObjectMeta, spec AuthSpec) Auth {
	return Auth{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[AuthSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "secrets.hashicorp.com/v1beta1",
				Kind:       "VaultAuth",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
