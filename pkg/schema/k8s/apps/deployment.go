package apps

import (
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type DeploymentSpec struct {
	Replicas int                  `yaml:",omitempty"`
	Selector meta.LabelSelector   `yaml:",omitempty"`
	Template core.PodTemplateSpec `yaml:",omitempty"`
}

type Deployment struct {
	shared.CommonK8sResourceWithSpec[DeploymentSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewDeployment(meta meta.ObjectMeta, spec DeploymentSpec) Deployment {
	return Deployment{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[DeploymentSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "apps/v1",
				Kind:       "Deployment",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
