package oci

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type RepoRef struct {
	Tag string
}

type RepoSpec struct {
	Url      string `yaml:",omitempty" validate:"required"`
	Interval string `yaml:",omitempty"`
	Ref      RepoRef
}

type Repo struct {
	shared.CommonK8sResourceWithSpec[RepoSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewRepo(meta meta.ObjectMeta, spec RepoSpec) Repo {
	return Repo{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[RepoSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "source.toolkit.fluxcd.io/v1beta2",
				Kind:       "OCIRepository",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
