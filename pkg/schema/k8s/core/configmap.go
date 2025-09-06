package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ConfigMap struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	Data                     map[string]string
}

func NewConfigMap(meta meta.ObjectMeta, data map[string]string) ConfigMap {
	return ConfigMap{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "v1",
			Kind:       "ConfigMap",
			Metadata:   meta,
		},
		Data: data,
	}
}
