package vaultsecretsoperator

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type AuthGlobalSpec struct {
	AllowedNamespaces []string   `yaml:"allowedNamespaces"`
	DefaultAuthMethod string     `yaml:"defaultAuthMethod"`
	Kubernetes        Kubernetes `yaml:"kubernetes"`
}

type AuthGlobal struct {
	shared.CommonK8sResourceWithSpec[AuthGlobalSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewAuthGlobal(meta meta.ObjectMeta, spec AuthGlobalSpec) AuthGlobal {
	return AuthGlobal{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[AuthGlobalSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "secrets.hashicorp.com/v1beta1",
				Kind:       "AuthGlobal",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
