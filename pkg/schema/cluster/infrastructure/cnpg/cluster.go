package cnpg

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type ClusterStorage struct {
	StorageClass string `yaml:"storageClass"`
	Size         string `yaml:"size"`
}

type InheritedMetadata struct {
	Annotations map[string]string `yaml:"annotations"`
}

type SuperuserSecret struct {
	Name string `yaml:"name"`
}

type InitDB struct {
	PostInitTemplateSQL []string `yaml:"postInitTemplateSQL"`
}

type Bootstrap struct {
	Initdb InitDB `yaml:"initdb"`
}

type AffinityConfiguration struct {
	NodeSelector map[string]string `yaml:"nodeSelector,omitempty"`
}

type ClusterSpec struct {
	Instances             int                   `yaml:"instances"`
	ImageName             string                `yaml:"imageName,omitempty"`
	Bootstrap             Bootstrap             `yaml:"bootstrap,omitempty"`
	Storage               ClusterStorage        `yaml:"storage"`
	InheritedMetadata     InheritedMetadata     `yaml:"inheritedMetadata,omitempty"`
	SuperuserSecret       SuperuserSecret       `yaml:"superuserSecret,omitempty"`
	EnableSuperuserAccess bool                  `yaml:"enableSuperuserAccess,omitempty"`
	Affinity              AffinityConfiguration `yaml:"affinity,omitempty"`
}

type Cluster struct {
	shared.CommonK8sResourceWithSpec[ClusterSpec] `yaml:",inline"`
}

func NewCluster(meta meta.ObjectMeta, spec ClusterSpec) Cluster {
	return Cluster{
		shared.CommonK8sResourceWithSpec[ClusterSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "postgresql.cnpg.io/v1",
				Kind:       "Cluster",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
