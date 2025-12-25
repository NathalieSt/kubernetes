package apps

import (
	"kubernetes/pkg/schema/k8s/core"
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type DaemonSetSpec struct {
	Selector meta.LabelSelector   `yaml:",omitempty"`
	Template core.PodTemplateSpec `yaml:",omitempty"`
}

type DaemonSet struct {
	shared.CommonK8sResourceWithSpec[DaemonSetSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewDaemonSet(meta meta.ObjectMeta, spec DaemonSetSpec) DaemonSet {
	return DaemonSet{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[DaemonSetSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "apps/v1",
				Kind:       "DaemonSet",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
