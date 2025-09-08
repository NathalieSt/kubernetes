package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ServiceAccount struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
}

func NewServiceAccount(meta meta.ObjectMeta) ServiceAccount {
	return ServiceAccount{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "v1",
			Kind:       "ServiceAccount",
			Metadata:   meta,
		},
	}
}
