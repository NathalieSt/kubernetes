package helm

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type RepoType string

const (
	Default RepoType = "default"
	OCI     RepoType = "oci"
)

type RepoSpec struct {
	RepoType RepoType `yaml:"type,omitempty"`
	Url      string   `yaml:",omitempty" validate:"required"`
	Interval string   `yaml:",omitempty"`
}

type Repo struct {
	shared.CommonK8sResourceWithSpec[RepoSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewRepo(meta meta.ObjectMeta, spec RepoSpec) Repo {
	return Repo{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[RepoSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "source.toolkit.fluxcd.io/v1",
				Kind:       "HelmRepository",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
