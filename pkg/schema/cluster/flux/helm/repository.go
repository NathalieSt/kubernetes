package helm

import "kubernetes/pkg/schema/k8s/meta"

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
	ApiVersion string          `yaml:"apiVersion" validate:"required"`
	Kind       string          `validate:"required"`
	Metadata   meta.ObjectMeta `validate:"required"`
	Spec       RepoSpec        `validate:"required"`
}

func NewRepo(meta meta.ObjectMeta, spec RepoSpec) Repo {
	return Repo{
		ApiVersion: "source.toolkit.fluxcd.io/v1",
		Kind:       "HelmRepository",
		Metadata:   meta,
		Spec:       spec,
	}
}
