package storage

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type StorageClassData struct {
	Provisioner          string            `yaml:"provisioner"`
	Parameters           map[string]string `yaml:"parameters"`
	ReclaimPolicy        string            `yaml:"reclaimPolicy"`
	VolumeBindingMode    string            `yaml:"volumeBindingMode"`
	AllowVolumeExpansion bool              `yaml:"allowVolumeExpansion"`
	MountOptions         []string          `yaml:"mountOptions"`
}

type StorageClass struct {
	shared.CommonK8sResource `yaml:",inline"`
	StorageClassData         `yaml:",inline"`
}

func NewStorageClass(meta meta.ObjectMeta, storageClassData StorageClassData) StorageClass {
	return StorageClass{
		CommonK8sResource: shared.CommonK8sResource{
			ApiVersion: "storage.k8s.io/v1",
			Kind:       "StorageClass",
			Metadata:   meta,
		},
		StorageClassData: storageClassData,
	}
}
