package core

import (
	"kubernetes/pkg/schema/k8s/meta"
	"kubernetes/pkg/schema/shared"
)

type VolumeResourceRequirements struct {
	Limits   map[string]string `yaml:",omitempty"`
	Requests map[string]string `yaml:",omitempty"`
}

type PersistentVolumeClaimSpec struct {
	AccessModes      []AccessModes              `yaml:"accessModes,omitempty"`
	Selector         meta.LabelSelector         `yaml:"selector,omitempty"`
	Resources        VolumeResourceRequirements `yaml:"resources,omitempty"`
	StorageClassName string                     `yaml:"storageClassName,omitempty"`
	VolumeName       string                     `yaml:"volumeName,omitempty"`
}

type PersistentVolumeClaim struct {
	shared.CommonK8sResourceWithSpec[PersistentVolumeClaimSpec] `yaml:",omitempty,inline" validate:"required"`
}

func NewPersistentVolumeClaim(meta meta.ObjectMeta, spec PersistentVolumeClaimSpec) PersistentVolumeClaim {
	return PersistentVolumeClaim{
		CommonK8sResourceWithSpec: shared.CommonK8sResourceWithSpec[PersistentVolumeClaimSpec]{
			CommonK8sResource: shared.CommonK8sResource{
				ApiVersion: "v1",
				Kind:       "PersistentVolumeClaim",
				Metadata:   meta,
			},
			Spec: spec,
		},
	}
}
