package apps

import (
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

type DeploymentSpec struct {
	Replicas int
	Selector meta.LabelSelector
	Template core.PodTemplateSpec
}

type Deployment struct {
	ApiVersion string `validate:"required"`
	Kind       string `validate:"required"`
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
