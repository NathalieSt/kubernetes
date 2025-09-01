package apps

import (
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
)

type DeploymentSpec struct {
	Replicas int
	Selector Selector
	Template core.PodTemplateSpec
}

type Deployment struct {
	ApiVersion string `validate:"required"`
	Kind       string `validate:"required"`
	Metadata   meta.ObjectMeta
	Spec       DeploymentSpec
}
