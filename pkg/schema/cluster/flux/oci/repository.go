package oci

import "kubernetes/pkg/schema/k8s/meta"

type RepoRef struct {
	Tag string
}

type RepoSpec struct {
	Url      string `yaml:",omitempty" validate:"required"`
	Interval string `yaml:",omitempty"`
	Ref      RepoRef
}

type Repo struct {
	ApiVersion string          `yaml:"apiVersion" validate:"required"`
	Kind       string          `validate:"required"`
	Metadata   meta.ObjectMeta `validate:"required"`
	Spec       RepoSpec        `validate:"required"`
}

func NewRepo(meta meta.ObjectMeta, spec RepoSpec) Repo {
	return Repo{
		ApiVersion: "source.toolkit.fluxcd.io/v1",
		Kind:       "OCIRepository",
		Metadata:   meta,
		Spec:       spec,
	}
}
