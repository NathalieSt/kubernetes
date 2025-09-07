package cnpg

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type DatabaseCluster struct {
	Name string `yaml:"name"`
}

type DatabaseSpec struct {
	Name    string          `yaml:"name"`
	Cluster DatabaseCluster `yaml:"cluster"`
	Owner   string          `yaml:"owner"`
}

type Database struct {
	shared.CommonK8sResourceWithSpec[DatabaseSpec] `yaml:",inline"`
}

func NewDatabase(meta meta.ObjectMeta, spec DatabaseSpec) Database {
	return Database{
		shared.CommonK8sResourceWithSpec[DatabaseSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "postgresql.cnpg.io/v1",
				Kind:       "Database",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
