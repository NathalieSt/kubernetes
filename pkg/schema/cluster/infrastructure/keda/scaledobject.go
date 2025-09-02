package keda

import "kubernetes/pkg/schema/k8s/meta"

type ScaleTargetRef struct {
	ApiVersion string `yaml:"apiVersion,omitempty"`
	Kind       string `yaml:",omitempty"`
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
	ApiVersion string `yaml:"apiVersion" validate:"required"`
	Kind       string
	Metadata   meta.ObjectMeta
	Spec       ScaledObjectSpec
}

func NewScaledObject(meta meta.ObjectMeta, spec ScaledObjectSpec) ScaledObject {
	return ScaledObject{
		ApiVersion: "keda.sh/v1alpha1",
		Kind:       "ScaledObject",
		Metadata:   meta,
		Spec:       spec,
	}
}
