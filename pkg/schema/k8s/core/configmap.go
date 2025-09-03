package core

import (
	"kubernetes/pkg/schema/shared"
)

type ConfigMap struct {
	shared.CommonK8sResource `yaml:",omitempty,inline" validate:"required"`
	Data                     map[string]string
}
