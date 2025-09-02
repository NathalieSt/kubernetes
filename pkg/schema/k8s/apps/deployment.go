package apps

import (
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

type DeploymentSpec struct {
	Replicas int                  `yaml:",omitempty"`
	Selector meta.LabelSelector   `yaml:",omitempty"`
	Template core.PodTemplateSpec `yaml:",omitempty"`
}

type Deployment struct {
	ApiVersion string `yaml:"apiVersion,omitempty" validate:"required"`
	Kind       string `yaml:",omitempty" validate:"required"`
	Metadata   meta.ObjectMeta
	Spec       DeploymentSpec
}

func NewDeployment(meta meta.ObjectMeta, spec DeploymentSpec) Deployment {
	return Deployment{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata:   meta,
		Spec:       spec,
	}
}
