package vaultsecretsoperator

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ConnectionSpec struct {
	Address string `yaml:"address,omitempty"`
}

type Connection struct {
	shared.CommonK8sResourceWithSpec[ConnectionSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewConnection(meta meta.ObjectMeta, spec ConnectionSpec) Connection {
	return Connection{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[ConnectionSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "secrets.hashicorp.com/v1beta1",
				Kind:       "VaultConnection",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
