package keda

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type Kind string

const (
	Deployment  Kind = "Deployment"
	StatefulSet Kind = "StatefulSet"
)

type ScaleTargetRef struct {
	ApiVersion string `yaml:"apiVersion,omitempty"`
	Kind       Kind   `yaml:",omitempty"`
	Name       string `yaml:",omitempty"`
}

type ScaledObjectTriggerMeta struct {
	Timezone        string
	Start           string
	End             string
	DesiredReplicas string `yaml:"desiredReplicas,omitempty"`
}

type KedaScaler string

const (
	Cron KedaScaler = "cron"
)

type ScaledObjectTrigger struct {
	ScalerType KedaScaler              `yaml:"type,omitempty" validate:"required"`
	Metadata   ScaledObjectTriggerMeta `yaml:",omitempty" validate:"required"`
}

type ScaledObjectSpec struct {
	ScaleTargetRef  ScaleTargetRef `yaml:"scaleTargetRef,omitempty"`
	MinReplicaCount int            `yaml:"minReplicaCount,"`
	CooldownPeriod  int            `yaml:"cooldownPeriod,omitempty"`
	Triggers        []ScaledObjectTrigger
}

type ScaledObject struct {
	shared.CommonK8sResourceWithSpec[ScaledObjectSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewScaledObject(meta meta.ObjectMeta, spec ScaledObjectSpec) ScaledObject {
	return ScaledObject{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[ScaledObjectSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "keda.sh/v1alpha1",
				Kind:       "ScaledObject",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
